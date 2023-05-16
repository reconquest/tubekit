package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/reconquest/karma-go"
	clientcmdapi "k8s.io/client-go/tools/clientcmd"
)

type Resource struct {
	Name      string
	Namespace string
}

func parseKubernetesContexts(kubeconfig string) ([]string, error) {
	loader := clientcmdapi.NewDefaultClientConfigLoadingRules()
	if kubeconfig != "" {
		loader.Precedence = append(
			[]string{kubeconfig},
			loader.Precedence...,
		)
	}

	config, err := loader.Load()
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to load kube config",
		)
	}

	contexts := []string{}
	for name := range config.Contexts {
		contexts = append(contexts, name)
	}

	return contexts, nil
}

func requestNamespaces(client string, params *Params) ([]string, error) {
	// omit namespace argument because requesting list of them
	cmd, args := getCommand(
		client,
		buildArgKubeconfig(params.Kubeconfig),
		buildArgContext(params.Context),
		"", "",
		"get", "namespaces", "-o", "json",
	)

	debugcmd(args)

	ctx := karma.Describe(
		"cmdline",
		fmt.Sprintf("%q", args),
	)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	contents, err := cmd.Output()
	if err != nil {
		return nil, ctx.Format(
			err,
			"tubectl command failed",
		)
	}

	resources, err := unmarshalResources(contents)
	if err != nil {
		return nil, ctx.Reason(err)
	}

	namespaces := []string{}
	for _, resource := range resources {
		namespaces = append(namespaces, resource.Name)
	}

	return namespaces, nil
}

func requestResources(client string, params *Params) ([]Resource, error) {
	cmd, args := getCommand(
		client,
		buildArgKubeconfig(params.Kubeconfig),
		buildArgContext(params.Context),
		buildArgNamespace(params.Namespace),
		buildArgAllNamespaces(params.AllNamespaces),
		"get", params.Match.Resource, "-o", "json",
	)

	debugcmd(args)

	ctx := karma.Describe(
		"cmdline",
		fmt.Sprintf("%q", args),
	)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	contents, err := cmd.Output()
	if err != nil {
		return nil, ctx.Format(
			err,
			"tubectl command failed",
		)
	}

	resources, err := unmarshalResources(contents)
	if err != nil {
		return nil, ctx.Reason(err)
	}

	return resources, nil
}

func unmarshalResources(contents []byte) ([]Resource, error) {
	var answer struct {
		Items []struct {
			Metadata Resource
		}
	}

	err := json.Unmarshal(contents, &answer)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to unmarshal JSON output",
		)
	}

	resources := []Resource{}
	for _, item := range answer.Items {
		resources = append(resources, item.Metadata)
	}

	return resources, nil
}

func getCommand(
	client string,
	argKubeconfig string,
	argContext string,
	argNamespace string,
	argAllNamespaces string,
	value ...string,
) (*exec.Cmd, []string) {
	args := []string{}
	if argKubeconfig != "" {
		args = append(args, argKubeconfig)
	}
	if argContext != "" {
		args = append(args, argContext)
	}
	if argNamespace != "" {
		args = append(args, argNamespace)
	}

	args = append(args, value...)

	if argAllNamespaces != "" {
		args = append(args, argAllNamespaces)
	}

	return exec.Command(client, args...), append([]string{client}, args...)
}

func buildArgKubeconfig(value string) string {
	if value != "" {
		return "--kubeconfig=" + value
	}

	return ""
}

func buildArgContext(value string) string {
	if value != "" {
		return "--context=" + value
	}

	return ""
}

func buildArgNamespace(value string) string {
	if value != "" {
		return "--namespace=" + value
	}

	return ""
}

func buildArgAllNamespaces(value bool) string {
	if value {
		return "--all-namespaces"
	}

	return ""
}
