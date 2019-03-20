package main

import (
	"fmt"
	"strings"
)

const (
	argPrefix = "--tube-"
)

type Flags struct {
	Help    bool
	Debug   bool
	Version bool
}

func parseFlags(raw []string) ([]string, *Flags, error) {
	args := make([]string, len(raw))
	copy(args, raw)

	flags := &Flags{}
	for i := 0; i < len(args); i++ {
		if i == 0 {
			continue
		}

		arg := args[i]

		if arg == "--" {
			break
		}

		if strings.HasPrefix(arg, argPrefix) {
			suffix := strings.TrimPrefix(arg, argPrefix)

			switch suffix {
			case "help":
				flags.Help = true

			case "debug":
				flags.Debug = true

			case "version":
				flags.Version = true

			default:
				return raw, nil, fmt.Errorf(
					"unexpected flag %q specified, see %shelp",
					arg, argPrefix,
				)
			}

			args = append(args[:i], args[i+1:]...)
			i--
		}
	}

	return args, flags, nil
}
