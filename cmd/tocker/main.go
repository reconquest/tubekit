package tocker

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/reconquest/executil-go"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/tubekit/internal/preflags"
)

var (
	version = "manual build"
	debug   = os.Getenv("TUBEKIT_DEBUG") == "1"
)

const (
	messageHelp = `tocker

tocker is a simple yet powerful wrapper around Docker which adds a bit of
magic to your everyday docker routines by reducing the complexity of working
with container names.

Usage:
	tocker [kubectl options]

Options:
  --tube-version  Show version of tocker.
  --tube-debug    Print debug messages.
  --tube-help     Show this message.

Docs: https://github.com/reconquest/tubekit`
)

func initFlags() {
	var flags *preflags.Flags
	var err error

	os.Args, flags, err = preflags.Parse(`--tube-`, os.Args)
	if err != nil {
		log.Fatalln(err)
	}

	switch {
	case flags.Help:
		fmt.Println(messageHelp)
		os.Exit(0)

	case flags.Debug:
		debug = true

	case flags.Version:
		fmt.Println(version)
		os.Exit(0)
	}
}

func getClientPath() (string, error) {
	env := os.Getenv("TUBEKIT_DOCKER")
	if env != "" {
		return env, nil
	}

	path, err := exec.LookPath("docker")
	if err != nil {
		return "", karma.Format(
			err,
			"unable to find docker in $PATH, have you installed it?",
		)
	}

	return path, nil
}

func main() {
	initFlags()

	client, err := getClientPath()
	if err != nil {
		log.Fatalln(err)
	}

	params := parseParams(os.Args)

	params, err = completeParams(client, params)
	if err != nil {
		log.Fatalln(err)
	}

	if params.Match == nil {
		syscallExec(client, params)
		return
	}

	resources, err := requestResources(client, params)
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

	tasks := getTasks(client, params, matched)

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

func syscallExec(client string, params *Params) {
	args := []string{client}

	if arg := buildArgKubeconfig(params.Kubeconfig); arg != "" {
		args = append(args, arg)
	}

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
		client,
		args,
		os.Environ(),
	)
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
