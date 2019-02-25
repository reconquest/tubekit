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

func requestNamespaces(ctlPath, context string) ([]string, error) {
	args := []string{ctlPath}
	if context != "" {
		args = append(args, "--context="+context)
	}

	args = append(
		args,
		"get", "namespaces", "-o", "json",
	)

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	contents, err := cmd.Output()
	if err != nil {
		return nil, karma.Format(
			err,
			"%q failed", args,
		)
	}

	var answer struct {
		Items []struct {
			Metadata struct {
				Name string `json:"name"`
			}
		}
	}

	err = json.Unmarshal(contents, &answer)
	if err != nil {
		return nil, karma.
			Describe(
				"cmdline",
				fmt.Sprintf("%q", args),
			).
			Format(
				err,
				"unable to unmarshal JSON output",
			)
	}

	namespaces := []string{}
	for _, item := range answer.Items {
		namespaces = append(namespaces, item.Metadata.Name)
	}

	return namespaces, nil
}

//func getCommand(context string, namespace string, values ...string) *exec.Cmd {
//}
