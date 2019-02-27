package main

import (
	"strconv"
	"strings"
	"unicode"
)

type Params struct {
	Args []string

	CompleteContext   bool
	Context           string
	CompleteNamespace bool
	Namespace         string
	AllNamespaces     bool

	Match *ParamsMatch
}

type ParamsMatch struct {
	Resource    string
	Query       string
	Parallel    bool
	Select      bool
	Element     int
	Placeholder int
}

const (
	SymbolContext     = '@'
	SymbolNamespace   = '+'
	SymbolMatch       = '%'
	SymbolMatchSelect = ':'
)

const (
	FlagContext             = "--context"
	FlagContextValue        = "--context="
	FlagNamespace           = "--namespace"
	FlagNamespaceValue      = "--namespace="
	FlagNamespaceShort      = "-n"
	FlagNamespaceShortValue = "-n="
	FlagAllNamespaces       = "--all-namespaces"
)

var (
	resourceMapping = map[string]string{
		"exec":         "pod",
		"logs":         "pod",
		"port-forward": "pod",
	}
)

func parseParams(raw []string) *Params {
	params := Params{}

	for index := 0; index < len(raw); index++ {
		// skip name of program, we don't need it since we are going to run
		// skube
		if index == 0 {
			continue
		}

		value := raw[index]
		if len(value) < 2 {
			// we are interested only in parsing something that starts
			// with special symbol and a symbol like @\w+ or it's -n
			params.Args = append(params.Args, value)
			continue
		}

		if params.Context == "" {
			var usedNext bool
			params.Context, usedNext = parseContext(value, raw, index)
			if usedNext {
				index++
			}

			if params.Context != "" {
				params.CompleteContext = value[0] == SymbolContext

				continue
			}
		}

		if params.Namespace == "" && !params.AllNamespaces {
			var usedNext bool
			params.Namespace, params.AllNamespaces, usedNext = parseNamespace(
				value, raw, index,
			)
			if usedNext {
				index++
			}

			if params.Namespace != "" {
				params.CompleteNamespace = value[0] == SymbolNamespace
				continue
			}

			if params.AllNamespaces {
				continue
			}
		}

		// index>1 because query can be passed after some entity like
		// get deployments nodejs%
		// if index == 1 then it's first argument after name of program
		if index > 1 && params.Match == nil {
			params.Match = parseMatch(value, params.Args)

			if params.Match != nil {
				continue
			}
		}

		params.Args = append(params.Args, value)
	}

	return &params
}

func parseMatch(
	value string,
	args []string,
) *ParamsMatch {
	var (
		placeholder = len(args)
	)

	if value[len(value)-1] == SymbolMatch {
		match := &ParamsMatch{
			Resource:    mapResource(args[len(args)-1]),
			Placeholder: placeholder,
		}

		if value[len(value)-2] == SymbolMatch {
			match.Parallel = true
			match.Query = value[:len(value)-2]
		} else {
			match.Query = value[:len(value)-1]
		}

		return match
	}

	var matchSelect bool
	var digitsStart int

	for i := len(value) - 1; i > 1; i-- {
		symbol := value[i]
		if symbol == SymbolMatchSelect {
			matchSelect = true
			continue
		}

		if matchSelect {
			if symbol != SymbolMatch {
				return nil
			}

			element, err := strconv.Atoi(value[digitsStart:])
			if err != nil {
				panic(
					"BUG: digits expected but got something weird: " +
						value[digitsStart:],
				)
			}

			// also need to cut % and :
			match := &ParamsMatch{
				Resource:    mapResource(args[len(args)-1]),
				Placeholder: placeholder,
				Query:       value[:digitsStart-2],
				Select:      true,
				Element:     element,
			}

			return match
		}

		if !unicode.IsDigit(rune(symbol)) {
			return nil
		}

		digitsStart = i
	}

	return nil
}

func mapResource(resource string) string {
	mapped, ok := resourceMapping[resource]
	if ok {
		return mapped
	}

	return resource
}

func parseContext(
	value string,
	args []string,
	index int,
) (name string, usedNext bool) {
	if value[0] == SymbolContext {
		return value[1:], false
	}

	if value == FlagContext {
		if index+1 <= len(args)-1 {
			return args[index+1], true
		}

		return "", false
	}

	if len(value) > len(FlagContextValue) &&
		strings.HasPrefix(value, FlagContextValue) {
		return value[len(FlagContextValue):], false
	}

	return "", false
}

func parseNamespace(
	value string,
	args []string,
	index int,
) (name string, all bool, usedNext bool) {
	if value[0] == SymbolNamespace {
		ns := value[1:]
		if ns == "+" {
			return "", true, false
		}
		return ns, false, false
	}

	if value == FlagAllNamespaces {
		return "", true, false
	}

	if value == FlagNamespaceShort {
		if index+1 <= len(args)-1 {
			return args[index+1], false, true
		}

		return "", false, false
	}

	if len(value) > len(FlagNamespaceShortValue) &&
		strings.HasPrefix(value, FlagNamespaceShortValue) {
		return value[len(FlagNamespaceShortValue):], false, false
	}

	if value == FlagNamespace {
		if index+1 <= len(args)-1 {
			return args[index+1], false, true
		}

		return "", false, false
	}

	if len(value) > len(FlagNamespaceValue) &&
		strings.HasPrefix(value, FlagNamespaceValue) {
		return value[len(FlagNamespaceValue):], false, false
	}

	return "", false, false
}
