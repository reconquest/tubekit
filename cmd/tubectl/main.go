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

var (
	debug = os.Getenv("TUBEKIT_DEBUG") == "1"
)

func main() {
	ctlPath := "/usr/bin/kubectl"
	if envPath := os.Getenv("TUBEKIT_KUBECTL"); envPath != "" {
		ctlPath = envPath
	}

	params := parseParams(os.Args)

	params, err := completeParams(ctlPath, params)
	if err != nil {
		log.Fatalln(err)
	}

	if params.Match == nil {
		syscallExec(ctlPath, params)
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

	tasks := getTasks(ctlPath, params, matched)

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

func matchResources(resources []Resource, params *ParamsMatch) ([]Resource, error) {
	exp, err := regexp.Compile(params.Query)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to compile regexp: %s", params.Query,
		)
	}

	matched := []Resource{}
	for _, resource := range resources {
		if exp.MatchString(resource.Name) {
			matched = append(matched, resource)
		}
	}

	return matched, nil
}

func syscallExec(ctlPath string, params *Params) {
	args := []string{ctlPath}

	if arg := buildArgContext(params.Context); arg != "" {
		args = append(args, arg)
	}

	if arg := buildArgNamespace(params.Namespace); arg != "" {
		args = append(args, arg)
	}

	args = append(args, params.Args...)

	if arg := buildArgAllNamespaces(params.AllNamespaces); arg != "" {
		args = append(args, arg)
	}

	debugcmd(args)

	syscall.Exec(
		ctlPath,
		args,
		os.Environ(),
	)
}

func completeParams(ctlPath string, params *Params) (*Params, error) {
	if params.CompleteContext {
		contexts, err := parseKubernetesContexts()
		if err != nil {
			return params, err
		}

		completed := complete(contexts, params.Context)
		if completed == "" && params.Context != "" {
			return params, fmt.Errorf(
				"unable to find such context: %s",
				params.Context,
			)
		}

		params.Context = completed
	}

	if params.CompleteNamespace {
		namespaces, err := requestNamespaces(ctlPath, params)
		if err != nil {
			return params, karma.Format(
				err,
				"unable to retrieve list of available namespaces",
			)
		}

		completed := complete(namespaces, params.Namespace)
		if completed == "" && params.Context != "" {
			return params, fmt.Errorf(
				"unable to find such namespace in context %s: %s",
				params.Context,
				params.Namespace,
			)
		}

		params.Namespace = completed
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

func debugcmd(args []string) {
	if debug {
		values := []string{}
		for _, arg := range args {
			values = append(values, fmt.Sprintf("%q", arg))
		}

		log.Printf("%s", strings.Join(values, " "))
	}
}
