package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"

	"github.com/reconquest/karma-go"
	clientcmdapi "k8s.io/client-go/tools/clientcmd"
)

func parseKubernetesContexts() ([]string, error) {
	config, err := clientcmdapi.NewDefaultClientConfigLoadingRules().Load()
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

	sort.Strings(contexts)

	return contexts, nil
}

func requestNamespaces(ctlPath string, params Params) ([]string, error) {
	// omit namespace argument because requesting list of them
	cmd, args := getCommand(
		ctlPath, buildArgContext(params), "",
		"get", "namespaces", "-o", "json",
	)

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
			"kubectl command failed",
		)
	}

	namespaces, err := unmarshalNames(contents)
	if err != nil {
		return nil, ctx.Reason(err)
	}

	return namespaces, nil
}

func requestResources(ctlPath string, params Params) ([]string, error) {
	cmd, args := getCommand(
		ctlPath,
		buildArgContext(params),
		buildArgNamespace(params),
		"get", params.Match.Resource, "-o", "json",
	)

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
			"kubectl command failed",
		)
	}

	resources, err := unmarshalNames(contents)
	if err != nil {
		return nil, ctx.Reason(err)
	}

	return resources, nil
}

func unmarshalNames(contents []byte) ([]string, error) {
	var answer struct {
		Items []struct {
			Metadata struct {
				Name string `json:"name"`
			}
		}
	}

	err := json.Unmarshal(contents, &answer)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to unmarshal JSON output",
		)
	}

	resources := []string{}
	for _, item := range answer.Items {
		resources = append(resources, item.Metadata.Name)
	}

	return resources, nil
}

func getCommand(
	ctlPath string,
	argContext,
	argNamespace string,
	value ...string,
) (*exec.Cmd, []string) {
	args := []string{}
	if argContext != "" {
		args = append(args, argContext)
	}
	if argNamespace != "" {
		args = append(args, argNamespace)
	}

	args = append(args, value...)

	return exec.Command(ctlPath, args...), append([]string{ctlPath}, args...)
}

func buildArgContext(params Params) string {
	if params.Context != "" {
		return "--context=" + params.Context
	}

	return ""
}

func buildArgNamespace(params Params) string {
	if params.AllNamespaces {
		return "--all-namespaces"
	} else if params.Namespace != "" {
		return "--namespace=" + params.Namespace
	}

	return ""
}
