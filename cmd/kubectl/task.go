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

func getTasks(ctlPath string, args []string, resources []string, placeholder int) []Task {
	tasks := []Task{}

	for i := range resources {
		resource := resources[i]

		task := func(writer io.Writer) error {
			values := append(
				[]string{},
				args[:placeholder-1]...,
			)

			values = append(values, resource)

			values = append(values, args[placeholder-1:]...)

			return run(ctlPath, values, writer)
		}

		tasks = append(tasks, task)
	}

	return tasks
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
