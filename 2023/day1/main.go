package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	_ "embed"
)

func main() {
	var inputPath, partFilter string
	flag.StringVar(&inputPath, "input", "input.txt", "intput file")
	flag.StringVar(&partFilter, "part", "", "part to run")
	flag.Parse()

	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		slog.Error("could not read file", "path", inputPath, "err", err)
		os.Exit(1)
	}

	for _, f := range []func(string) (any, string, error){Part1, Part2} {
		funcName := strings.Split(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), ".")[1]
		if !strings.HasSuffix(funcName, partFilter) {
			continue
		}

		start := time.Now()
		result, debug, err := f(string(inputData))
		duration := time.Since(start)

		if err != nil {
			slog.Error("could not run part", "func", funcName, "err", err)
			break
		}

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

		if debug != "" {
			debugFile, err := os.OpenFile(fmt.Sprintf("debug-%s.txt", funcName), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0777)
			if err != nil {
				slog.Error("could not create or append debug file", "err", err)
				os.Exit(1)
			}
			defer debugFile.Close()

			_, err = debugFile.WriteString(debug)
			if err != nil {
				slog.Error("could not write debug", "func", funcName, "err", err)
				break
			}
		}

		slog.Info("finished running part", "func", funcName, "duration", duration, "result", result)
	}
}

func iterLines(input string) func(func(int, string) bool) {
	return func(yield func(int, string) bool) {
		lines := strings.Split(input, "\n")
		for i := 0; i < len(lines); i++ {
			line := strings.TrimSpace(lines[i])
			if line == "" {
				continue
			}

			if !yield(i, line) {
				return
			}
		}
	}
}

func Part1(input string) (any, string, error) {
	result := 0
	debug := ""

	for i, line := range iterLines(input) {
		var first, last string
		for _, r := range line {
			c := string(r)
			_, err := strconv.Atoi(c)
			if err != nil {
				continue
			}

			if first == "" {
				first = c
				continue
			}

			last = c
		}

		if last == "" {
			last = first
		}

		num, err := strconv.Atoi(first + last)
		if err != nil {
			return result, debug, fmt.Errorf("could not convert line %d %#v to number: %w", i, line, err)
		}

		result += num
		debug += fmt.Sprintf("%s => +%d | %d\n", line, num, result)
	}

	return result, debug, nil
}

func Part2(input string) (any, string, error) {
	result := 0
	debug := ""
	replacements := map[string]string{
		"one":   "1",
		"two":   "2",
		"three": "3",
		"four":  "4",
		"five":  "5",
		"six":   "6",
		"seven": "7",
		"eight": "8",
		"nine":  "9",
	}

	for i, line := range iterLines(input) {
		var first, last, past string
		setNumbers := func(c string) {
			if first == "" {
				first = c
				return
			}

			last = c
		}

		for _, r := range line {
			c := string(r)
			_, err := strconv.Atoi(c)
			if err != nil {
				past += c

				for old, new := range replacements {
					if strings.HasSuffix(past, old) {
						setNumbers(new)
						break
					}
				}

				continue
			}

			past = ""
			setNumbers(c)

		}

		if last == "" {
			last = first
		}

		num, err := strconv.Atoi(first + last)
		if err != nil {
			return result, debug, fmt.Errorf("could not convert line %d %#v to number: %w", i, line, err)
		}

		result += num
		debug += fmt.Sprintf("%s => +%d | %d\n", line, num, result)
	}

	return result, debug, nil
}
