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
				Context: "ctx",
				Args:    []string{"get", "pods"},
			},
		},
		{
			[]string{"@ctx", "get", "pods"},
			Params{
				Context: "ctx",
				Args:    []string{"get", "pods"},
			},
		},
		{
			[]string{"get", "@ctx", "pods"},
			Params{
				Context: "ctx",
				Args:    []string{"get", "pods"},
			},
		},
		{
			[]string{"@ctx"},
			Params{
				Context: "ctx",
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
				Namespace: "ns",
				Args:      []string{"get", "pods"},
			},
		},
		{
			[]string{"get", "+ns", "pods"},
			Params{
				Namespace: "ns",
				Args:      []string{"get", "pods"},
			},
		},
		{
			[]string{"+ns", "get", "pods"},
			Params{
				Namespace: "ns",
				Args:      []string{"get", "pods"},
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
				Namespace: "ns",
				Args:      []string{"get", "pods", "++"},
			},
		},
		// weird case, no idea what behaviour could be expected in such case
		{
			[]string{"get", "pods", "+ns", "--all-namespaces"},
			Params{
				Namespace: "ns",
				Args:      []string{"get", "pods", "--all-namespaces"},
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
				Context:   "ctx",
				Namespace: "ns",
				Args:      []string{"get", "pods", "-v"},
			},
		},
	}

	for _, testcase := range testcases {
		params := parseParams(append([]string{"binary"}, testcase.raw...))
		if !test.EqualValues(testcase.result, params, "%q", testcase.raw) {
			t.FailNow()
		}
	}
}
