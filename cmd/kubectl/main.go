package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"syscall"

	"github.com/reconquest/executil-go"
	"github.com/reconquest/karma-go"
)

func main() {
	ctlPath := "/usr/bin/kubectl"
	if envPath := os.Getenv("QUADRO_KUBECTL"); envPath != "" {
		ctlPath = envPath
	}

	params := parseParams(os.Args)

	params, err := completeParams(ctlPath, params)
	if err != nil {
		log.Fatalln(err)
	}

	args := buildArgs(params)

	if params.Match == nil {
		syscallExec(ctlPath, args)
		return
	}

	resources, err := requestResources(ctlPath, params)
	if err != nil {
		log.Fatalln(err)
	}

	matched, err := matchResources(resources, params.Match)
	if err != nil {
		log.Fatalln(err)
	}

	if params.Match.Select {
		if params.Match.Element < 1 || params.Match.Element > len(matched) {
			log.Fatalf(
				"no resource with such index: %d, "+
					"matched resources: %q (total %d)",
				params.Match.Element,
				matched,
				len(matched),
			)
		}

		matched = matched[params.Match.Element-1 : params.Match.Element]
	}

	if len(matched) == 0 {
		log.Fatalf(
			"no resources found: %s %s",
			params.Match.Resource,
			params.Match.Query,
		)
	}

	tasks := getTasks(ctlPath, args, matched, params.Match.Placeholder)

	if params.Match.Parallel {
		parallelize(tasks)
		return
	}

	code := 0
	for _, task := range tasks {
		err := task(os.Stdout)
		if err != nil {
			log.Println(err)

			if executil.IsExitError(err) {
				code = executil.GetExitStatus(err)
			}
		}
	}

	os.Exit(code)
}

func matchResources(resources []string, params *ParamsMatch) ([]string, error) {
	exp, err := regexp.Compile(params.Query)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to compile regexp: %s", params.Query,
		)
	}

	matched := []string{}
	for _, resource := range resources {
		if exp.MatchString(resource) {
			matched = append(matched, resource)
		}
	}

	return matched, nil
}

func syscallExec(ctlPath string, args []string) {
	syscall.Exec(
		ctlPath,
		append([]string{"kubectl"}, args...),
		os.Environ(),
	)
}

func buildArgs(params Params) []string {
	args := []string{}

	if arg := buildArgContext(params); arg != "" {
		args = append(args, arg)
	}

	if arg := buildArgNamespace(params); arg != "" {
		args = append(args, arg)
	}

	args = append(args, params.Args...)

	return args
}

func completeParams(ctlPath string, params Params) (Params, error) {
	if params.CompleteContext {
		contexts, err := parseKubernetesContexts()
		if err != nil {
			return params, err
		}

		completeed := complete(contexts, params.Context)
		if completeed == "" && params.Context != "" {
			return params, fmt.Errorf(
				"unable to find such context: %s",
				params.Context,
			)
		}

		params.Context = completeed
	}

	if params.CompleteNamespace {
		namespaces, err := requestNamespaces(ctlPath, params)
		if err != nil {
			return params, karma.Format(
				err,
				"unable to retrieve list of available namespaces",
			)
		}

		completeed := complete(namespaces, params.Namespace)
		if completeed == "" && params.Context != "" {
			return params, fmt.Errorf(
				"unable to find such namespace in context %s: %s",
				params.Context,
				params.Namespace,
			)
		}

		params.Namespace = completeed
	}

	return params, nil
}

func complete(items []string, query string) string {
	var partial string
	for _, item := range items {
		if item == query {
			return item
		}

		if strings.Contains(item, query) {
			partial = item
		}
	}

	return partial
}
