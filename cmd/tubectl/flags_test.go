package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlags(t *testing.T) {
	test := assert.New(t)

	testcases := []struct {
		raw   []string
		after []string
		flags *Flags
	}{
		{
			[]string{"x", "--version"},
			[]string{"x", "--version"},
			&Flags{},
		},
		{
			[]string{"x", "a", "b", "c%%", "@d", "--", "e", "", "v"},
			[]string{"x", "a", "b", "c%%", "@d", "--", "e", "", "v"},
			&Flags{},
		},
		{
			[]string{"x", "--version", "--tube-version", "@after"},
			[]string{"x", "--version", "@after"},
			&Flags{
				Version: true,
			},
		},
		{
			[]string{"x", "--version", "--tube-debug", "@after"},
			[]string{"x", "--version", "@after"},
			&Flags{
				Debug: true,
			},
		},
		{
			[]string{"x", "--version", "--tube-help", "@after"},
			[]string{"x", "--version", "@after"},
			&Flags{
				Help: true,
			},
		},
		{
			[]string{"x", "--version", "--tube-help"},
			[]string{"x", "--version"},
			&Flags{
				Help: true,
			},
		},
		{
			[]string{"x", "--version", "--tube-help", "--tube-version", "--tube-debug"},
			[]string{"x", "--version"},
			&Flags{
				Help:    true,
				Version: true,
				Debug:   true,
			},
		},
		{
			[]string{"x", "--version", "--", "--tube-help"},
			[]string{"x", "--version", "--", "--tube-help"},
			&Flags{},
		},
		{
			[]string{"x", "--tube-unexpected"},
			[]string{"x", "--tube-unexpected"},
			nil,
		},
	}

	for _, testcase := range testcases {
		args, flags, err := parseFlags(testcase.raw)
		if !test.EqualValues(testcase.flags, flags, "%q", testcase.raw) {
			t.FailNow()
		}

		if !test.EqualValues(testcase.after, args, "%q", testcase.raw) {
			t.FailNow()
		}

		if testcase.flags == nil {
			if !test.Error(err) {
				t.FailNow()
			}
		}
	}
}
