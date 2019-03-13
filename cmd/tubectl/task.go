package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/reconquest/karma-go"
	"github.com/reconquest/lineflushwriter-go"
)

type Task func(io.Writer) error

func getTasks(
	ctlPath string,
	params *Params,
	resources []Resource,
) []Task {
	tasks := []Task{}

	for _, resource := range resources {
		tasks = append(
			tasks,
			getTask(ctlPath, params, resource),
		)
	}

	return tasks
}

func getTask(ctlPath string, params *Params, resource Resource) Task {
	return func(writer io.Writer) error {
		values := []string{}

		if arg := buildArgContext(params.Context); arg != "" {
			values = append(values, arg)
		}

		if arg := buildArgNamespace(resource.Namespace); arg != "" {
			values = append(values, arg)
		}

		values = append(
			values,
			params.Args[:params.Match.Placeholder]...,
		)

		values = append(values, resource.Name)

		values = append(values, params.Args[params.Match.Placeholder:]...)

		return run(ctlPath, values, writer)
	}
}

func parallelize(tasks []Task) {
	workers := &sync.WaitGroup{}
	workers.Add(len(tasks))
	for _, task := range tasks {
		writer := lineflushwriter.New(
			os.Stdout,
			&sync.Mutex{},
			true,
		)

		go func(task Task) {
			defer workers.Done()

			task(writer)
		}(task)
	}

	workers.Wait()
}

func run(ctlPath string, args []string, writer io.Writer) error {
	debugcmd(append([]string{ctlPath}, args...))

	cmd := exec.Command(ctlPath, args...)
	cmd.Stdout = writer
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return karma.
			Describe(
				"cmdline",
				fmt.Sprintf("%q", append([]string{ctlPath}, args...)),
			).Format(err, "command failed")
	}

	return nil
}
