package main

import (
	"encoding/csv"
	"strings"
)

type Args []string

func ParseArgs(input string) Args {
	r := csv.NewReader(strings.NewReader(input))
	s, err := r.Read()
	if err != nil {
		return strings.Fields(input)
	}

	return s
}

// Name returns the command name, $0
func (a Args) Name() string {
	if len(a) == 0 {
		return ""
	}

	return a[0]
}

// Args returns all the arguments, $@
func (a Args) Args() []string {
	if len(a) <= 1 {
		return []string{}
	}

	return a[1:]
}
