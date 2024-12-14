package main

import (
	"flag"
	"strings"
)

func parseArgs() Args {
	var args Args
	flag.StringVar(&args.Input, "input", "input.txt", "")
	flag.StringVar(&args.PartFilter, "part", "", "")
	flag.BoolVar(&args.Verbose, "v", false, "verbose debug")
	flag.Parse()

	// args.Input = "input-example-1.txt"
	// args.Part = "1"

	return args
}

func Part1(input string, debug *Debugger) (any, error) {
	result := 0

	for _, line := range strings.Split(input, "\n") {
		_ = line
	}

	return result, nil
}

func Part2(input string, debug *Debugger) (any, error) {
	result := 0

	for _, line := range strings.Split(input, "\n") {
		_ = line
	}

	return result, nil
}
