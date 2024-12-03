package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	_ "embed"
)

func main() {
	inputPath := flag.String("input", "input.txt", "intput file")
	flag.Parse()

	inputData, err := os.ReadFile(*inputPath)
	if err != nil {
		slog.Error("could not read file", "path", *inputPath, "err", err)
		os.Exit(1)
	}

	for _, f := range []func(string) (any, string, error){Part1, Part2} {
		funcName := strings.Split(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), ".")[1]

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

		slog.Info("finished running part", "func", funcName, "duration", duration, "result", result)
	}
}

func Part1(input string) (any, string, error) {
	re := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	result := 0
	debug := ""

	for _, match := range re.FindAllStringSubmatch(input, -1) {
		num1, _ := strconv.Atoi(match[1])
		num2, _ := strconv.Atoi(match[2])
		subResult := num1 * num2
		result += subResult
		debug += fmt.Sprintf("%d * %d = %d | %d | %v\n", num1, num2, subResult, result, match)
	}

	return result, debug, nil
}

func Part2(input string) (any, string, error) {
	re := regexp.MustCompile(`do\(\)|don't\(\)|mul\((\d{1,3}),(\d{1,3})\)`)
	result := 0
	debug := ""

	mulEnabled := true
	for _, match := range re.FindAllStringSubmatch(input, -1) {
		debug += match[0]
		if match[0] == "do()" {
			mulEnabled = true
		} else if match[0] == "don't()" {
			mulEnabled = false
		} else if mulEnabled {
			num1, _ := strconv.Atoi(match[1])
			num2, _ := strconv.Atoi(match[2])
			subResult := num1 * num2
			result += subResult
			debug += fmt.Sprintf(" | %d * %d = %d | %d", num1, num2, subResult, result)
		} else {
			debug += " | SKIP"
		}
		debug += "\n"
	}

	return result, debug, nil
}
