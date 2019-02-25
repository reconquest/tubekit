package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParams_Suites(t *testing.T) {
	test := assert.New(t)

	testcases := []struct {
		raw    []string
		result Params
	}{
		{
			[]string{"get", "pods", "@ctx"},
			Params{
				CompleteContext: true,
				Context:         "ctx",
				Args:            []string{"get", "pods"},
			},
		},
		{
			[]string{"@ctx", "get", "pods"},
			Params{
				CompleteContext: true,
				Context:         "ctx",
				Args:            []string{"get", "pods"},
			},
		},
		{
			[]string{"get", "@ctx", "pods"},
			Params{
				CompleteContext: true,
				Context:         "ctx",
				Args:            []string{"get", "pods"},
			},
		},
		{
			[]string{"@ctx"},
			Params{
				CompleteContext: true,
				Context:         "ctx",
			},
		},
		{
			[]string{"get", "pods"},
			Params{
				Args: []string{"get", "pods"},
			},
		},
		{
			[]string{"get", "pods", "@"},
			Params{
				Args: []string{"get", "pods", "@"},
			},
		},
		{
			[]string{"get", "pods", "-c", "ctx"},
			Params{
				Args: []string{"get", "pods", "-c", "ctx"},
			},
		},
		{
			[]string{"get", "pods", "-c=ctx"},
			Params{
				Args: []string{"get", "pods", "-c=ctx"},
			},
		},
		{
			[]string{"get", "pods", "--context", "ctx"},
			Params{
				Context: "ctx",
				Args:    []string{"get", "pods"},
			},
		},
		{
			[]string{"get", "pods", "--context"},
			Params{
				Args: []string{"get", "pods", "--context"},
			},
		},
		{
			[]string{"get", "pods", "--context=ctx"},
			Params{
				Context: "ctx",
				Args:    []string{"get", "pods"},
			},
		},
		{
			[]string{"get", "pods", "--context=", "ctx"},
			Params{
				Args: []string{"get", "pods", "--context=", "ctx"},
			},
		},
		{
			[]string{"get", "pods", "+ns"},
			Params{
				CompleteNamespace: true,
				Namespace:         "ns",
				Args:              []string{"get", "pods"},
			},
		},
		{
			[]string{"get", "+ns", "pods"},
			Params{
				CompleteNamespace: true,
				Namespace:         "ns",
				Args:              []string{"get", "pods"},
			},
		},
		{
			[]string{"+ns", "get", "pods"},
			Params{
				CompleteNamespace: true,
				Namespace:         "ns",
				Args:              []string{"get", "pods"},
			},
		},
		{
			[]string{"get", "pods", "-n", "ns"},
			Params{
				Namespace: "ns",
				Args:      []string{"get", "pods"},
			},
		},
		{
			[]string{"get", "pods", "--namespace", "ns"},
			Params{
				Namespace: "ns",
				Args:      []string{"get", "pods"},
			},
		},
		{
			[]string{"get", "pods", "--namespace=ns"},
			Params{
				Namespace: "ns",
				Args:      []string{"get", "pods"},
			},
		},
		{
			[]string{"get", "pods", "--namespace="},
			Params{
				Args: []string{"get", "pods", "--namespace="},
			},
		},
		{
			[]string{"get", "pods", "++"},
			Params{
				AllNamespaces: true,
				Args:          []string{"get", "pods"},
			},
		},
		{
			[]string{"++", "get", "pods"},
			Params{
				AllNamespaces: true,
				Args:          []string{"get", "pods"},
			},
		},
		{
			[]string{"++", "get", "pods", "+ns"},
			Params{
				AllNamespaces: true,
				Args:          []string{"get", "pods", "+ns"},
			},
		},
		{
			[]string{"get", "pods", "+ns", "++"},
			Params{
				CompleteNamespace: true,
				Namespace:         "ns",
				Args:              []string{"get", "pods", "++"},
			},
		},
		// weird case, no idea what behaviour could be expected in such case
		{
			[]string{"get", "pods", "+ns", "--all-namespaces"},
			Params{
				CompleteNamespace: true,
				Namespace:         "ns",
				Args:              []string{"get", "pods", "--all-namespaces"},
			},
		},
		{
			[]string{"--all-namespaces", "get", "pods", "+ns"},
			Params{
				AllNamespaces: true,
				Args:          []string{"get", "pods", "+ns"},
			},
		},

		//
		{
			[]string{"get", "pods", "-n", "ns", "@ctx", "-v"},
			Params{
				CompleteContext: true,
				Context:         "ctx",
				Namespace:       "ns",
				Args:            []string{"get", "pods", "-v"},
			},
		},
		{
			[]string{"get", "pods", "+ns", "@ctx", "-v"},
			Params{
				CompleteContext:   true,
				Context:           "ctx",
				CompleteNamespace: true,
				Namespace:         "ns",
				Args:              []string{"get", "pods", "-v"},
			},
		},
		//
		{
			[]string{"describe", "pods", "qu%"},
			Params{
				Match: &ParamsMatch{
					Query:  "qu",
					Entity: "pods",
				},
				Args: []string{"describe", "pods"},
			},
		},
		{
			[]string{"describe", "pods", "qu%%"},
			Params{
				Match: &ParamsMatch{
					Query:    "qu",
					Entity:   "pods",
					Parallel: true,
				},
				Args: []string{"describe", "pods"},
			},
		},
		{
			[]string{"describe", "pods", "qu%:1"},
			Params{
				Match: &ParamsMatch{
					Query:   "qu",
					Entity:  "pods",
					Select:  true,
					Element: 1,
				},
				Args: []string{"describe", "pods"},
			},
		},
		{
			[]string{"describe", "pods", "qu%:10"},
			Params{
				Match: &ParamsMatch{
					Query:   "qu",
					Entity:  "pods",
					Select:  true,
					Element: 10,
				},
				Args: []string{"describe", "pods"},
			},
		},
		{
			[]string{"describe", "pods", "qu2%:10"},
			Params{
				Match: &ParamsMatch{
					Query:   "qu2",
					Entity:  "pods",
					Select:  true,
					Element: 10,
				},
				Args: []string{"describe", "pods"},
			},
		},
		{
			[]string{"describe", "pods", "qu%10"},
			Params{
				Args: []string{"describe", "pods", "qu%10"},
			},
		},
		{
			[]string{"describe", "pods", "qu:10"},
			Params{
				Args: []string{"describe", "pods", "qu:10"},
			},
		},
		//
	}

	for _, testcase := range testcases {
		params := parseParams(append([]string{"binary"}, testcase.raw...))
		if !test.EqualValues(testcase.result, params, "%q", testcase.raw) {
			t.FailNow()
		}
	}
}
