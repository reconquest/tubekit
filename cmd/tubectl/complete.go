package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/reconquest/karma-go"
)

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

func completeParams(client string, params *Params) (*Params, error) {
	if params.CompleteContext {
		contexts, err := parseKubernetesContexts(params.Kubeconfig)
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
		namespaces, err := requestNamespaces(client, params)
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
