package main

import (
	"strings"
)

type Params struct {
	Context       string
	Namespace     string
	AllNamespaces bool
	Query         string
	Args          []string
}

const (
	flagContext             = "--context"
	flagContextValue        = "--context="
	flagNamespace           = "--namespace"
	flagNamespaceValue      = "--namespace="
	flagNamespaceShort      = "-n"
	flagNamespaceShortValue = "-n="
	flagAllNamespaces       = "--all-namespaces"
)

func parseParams(raw []string) Params {
	params := Params{}

	pushArg := func(value string) {
		params.Args = append(params.Args, value)
	}

	for i := 0; i < len(raw); i++ {
		// skip name of program, we don't need it since we are going to run
		// kubectl
		if i == 0 {
			continue
		}

		value := raw[i]
		if len(value) < 2 {
			// we are interested only in parsing something that starts
			// with special symbol and a symbol like @\w+ or it's -n
			pushArg(value)
			continue
		}

		if params.Context == "" {
			var usedNext bool
			params.Context, usedNext = parseContext(value, raw, i)
			if usedNext {
				i++
			}

			if params.Context != "" {
				continue
			}
		}

		if params.Namespace == "" && !params.AllNamespaces {
			var usedNext bool
			params.Namespace, params.AllNamespaces, usedNext = parseNamespace(value, raw, i)
			if usedNext {
				i++
			}

			if params.Namespace != "" || params.AllNamespaces {
				continue
			}
		}

		params.Args = append(params.Args, value)
	}

	return params
}

func parseContext(
	value string,
	raw []string,
	i int,
) (name string, usedNext bool) {
	if value[0] == '@' {
		return value[1:], false
	}

	if value == flagContext {
		if i+1 <= len(raw)-1 {
			return raw[i+1], true
		}

		return "", false
	}

	if len(value) > len(flagContextValue) &&
		strings.HasPrefix(value, flagContextValue) {
		return value[len(flagContextValue):], false
	}

	return "", false
}

func parseNamespace(
	value string,
	raw []string,
	i int,
) (name string, all bool, usedNext bool) {
	if value[0] == '+' {
		ns := value[1:]
		if ns == "+" {
			return "", true, false
		}
		return ns, false, false
	}

	if value == flagAllNamespaces {
		return "", true, false
	}

	if value == flagNamespaceShort {
		if i+1 <= len(raw)-1 {
			return raw[i+1], false, true
		}

		return "", false, false
	}

	if len(value) > len(flagNamespaceShortValue) &&
		strings.HasPrefix(value, flagNamespaceShortValue) {
		return value[len(flagNamespaceShortValue):], false, false
	}

	if value == flagNamespace {
		if i+1 <= len(raw)-1 {
			return raw[i+1], false, true
		}

		return "", false, false
	}

	if len(value) > len(flagNamespaceValue) &&
		strings.HasPrefix(value, flagNamespaceValue) {
		return value[len(flagNamespaceValue):], false, false
	}

	return "", false, false
}
