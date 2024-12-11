package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func main() {
	defaultInput := "input.txt"
	defaultPart := ""

	// defaultInput = "input-example-1.txt"
	// defaultPart = "1"

	var verboseDebug bool
	var inputPath, partFilter string
	flag.StringVar(&inputPath, "input", defaultInput, "")
	flag.StringVar(&partFilter, "part", defaultPart, "")
	flag.BoolVar(&verboseDebug, "v", false, "verbose debug")
	flag.Parse()

	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		slog.Error("could not read file", "path", inputPath, "err", err)
		os.Exit(1)
	}

	input := strings.TrimSuffix(strings.ReplaceAll(string(inputData), "\r\n", "\n"), "\n")

	for _, f := range []func(string, *Debugger) (any, error){Part1, Part2} {
		funcName := strings.Split(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), ".")[1]
		if !strings.HasSuffix(funcName, partFilter) {
			continue
		}

		debug := NewDebugBuilder(verboseDebug, fmt.Sprintf("debug-%s.txt", funcName), -1)
		defer debug.Close()

		start := time.Now()
		result, err := f(input, debug)
		duration := time.Since(start)

		if err != nil {
			slog.Error("could not run part", "func", funcName, "err", err)
			break
		}

		debug.Close()
		slog.Info("finished running part", "func", funcName, "duration", duration, "result", result)
	}
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
