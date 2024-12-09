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

	_ "embed"
)

const VisualizeStep = "==========STEP==========\n"

func main() {
	defaultInput := "input.txt"
	defaultPart := ""

	// defaultInput = "input-example-1.txt"
	// defaultPart = "1"

	var verboseDebug, outputFile bool
	var inputPath, partFilter string
	flag.StringVar(&inputPath, "input", defaultInput, "")
	flag.StringVar(&partFilter, "part", defaultPart, "")
	flag.BoolVar(&verboseDebug, "v", false, "verbose debug")
	flag.BoolVar(&outputFile, "o", false, "write output file")
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

		if outputFile {
			file, err := os.OpenFile(fmt.Sprintf("output-%s.txt", funcName), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
			if err != nil {
				slog.Error("could not create or append output file", "err", err)
				os.Exit(1)
			}
			defer file.Close()

			_, err = file.WriteString(fmt.Sprintf("%v", result))
			if err != nil {
				slog.Error("could not write result", "func", funcName, "result", result, "err", err)
				break
			}
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
