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
	flag.StringVar(&inputPath, "input", "input.txt", "")
	flag.StringVar(&partFilter, "part", "", "")
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

	maxes := map[string]int{
		"red":   12,
		"green": 13,
		"blue":  14,
	}

	for gameNum, line := range iterLines(input) {
		gameNum++
		success := true
		for _, set := range strings.Split(strings.Split(line, ":")[1], ";") {
			for _, cube := range strings.Split(set, ",") {
				parts := strings.Split(strings.TrimSpace(cube), " ")
				num, _ := strconv.Atoi(parts[0])
				if num > maxes[parts[1]] {
					success = false
					break
				}
			}
			if !success {
				break
			}
		}

		if success {
			result += gameNum
		}

		debug += fmt.Sprintf("%d %t\n", gameNum, success)
	}

	return result, debug, nil
}

func Part2(input string) (any, string, error) {
	result := 0
	debug := ""

	for gameNum, line := range iterLines(input) {
		gameNum += 1
		maxes := map[string]int{
			"red":   0,
			"green": 0,
			"blue":  0,
		}
		for _, set := range strings.Split(strings.Split(line, ":")[1], ";") {
			for _, cube := range strings.Split(set, ",") {
				parts := strings.Split(strings.TrimSpace(cube), " ")
				num, _ := strconv.Atoi(parts[0])
				if num > maxes[parts[1]] {
					maxes[parts[1]] = num
				}
			}
		}

		gameResult := 1
		for _, num := range maxes {
			gameResult *= num
		}

		result += gameResult
		debug += fmt.Sprintf("%d +%d => %d\n", gameNum, gameResult, result)
	}

	return result, debug, nil
}
