package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

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

	run(ctlPath, args)
}

func run(ctlPath string, args []string) {
	syscall.Exec(
		ctlPath,
		append([]string{"kubectl"}, args...),
		os.Environ(),
	)
}

func buildArgs(params Params) []string {
	args := []string{}

	if params.Context != "" {
		args = append(args, "--context="+params.Context)
	}

	if params.AllNamespaces {
		args = append(args, "--all-namespaces")
	} else if params.Namespace != "" {
		args = append(args, "--namespace="+params.Namespace)
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
		namespaces, err := requestNamespaces(ctlPath, params.Context)
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
