package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"runtime"
	"slices"
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
		lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
		for i := 0; i < len(lines); i++ {
			if !yield(i, lines[i]) {
				return
			}
		}
	}
}

func Part1(input string) (any, string, error) {
	result := 0
	debug := ""

	inRules := true
	pageRules := map[string][]string{}
	for _, line := range iterLines(input) {
		if line == "" {
			inRules = false
			for k, v := range pageRules {
				debug += fmt.Sprintf("%v: %v\n", k, v)
			}
			debug += "\n"
			continue
		}

		if inRules {
			parts := strings.Split(line, "|")
			pageRules[parts[1]] = append(pageRules[parts[1]], parts[0])
		} else {
			valid := true
			pages := strings.Split(line, ",")
			for i, page := range pages {
				rules := pageRules[page]
				for _, rule := range rules {
					if slices.Contains(pages[i+1:], rule) {
						valid = false
						break
					}
				}

				if !valid {
					break
				}
			}

			debug += fmt.Sprintf("%v", pages)
			if valid {
				middle, _ := strconv.Atoi(pages[len(pages)/2])
				result += middle

				debug += fmt.Sprintf(" +%d => %d", middle, result)
			}
			debug += "\n"
		}
	}

	return result, debug, nil
}

func Part2(input string) (any, string, error) {
	result := 0
	debug := ""

	inRules := true
	pageRules := map[string][]string{}
	var failedUpdates [][]string
	for _, line := range iterLines(input) {
		if line == "" {
			inRules = false
			for k, v := range pageRules {
				debug += fmt.Sprintf("%v: %v\n", k, v)
			}
			debug += "\n"
			continue
		}

		if inRules {
			parts := strings.Split(line, "|")
			pageRules[parts[1]] = append(pageRules[parts[1]], parts[0])
		} else {
			valid := true
			update := strings.Split(line, ",")
			for i, page := range update {
				rules := pageRules[page]
				for _, rule := range rules {
					if slices.Contains(update[i+1:], rule) {
						valid = false
						break
					}
				}

				if !valid {
					break
				}
			}

			debug += fmt.Sprintf("%v %t\n", update, valid)
			if !valid {
				failedUpdates = append(failedUpdates, update)
			}
		}
	}

	debug += "\n"
	for _, update := range failedUpdates {
		newUpdate := pt2FixOrder(update, pageRules)

		middle, _ := strconv.Atoi(newUpdate[len(newUpdate)/2])
		result += middle

		debug += fmt.Sprintf("%v => %v +%d = %d\n", update, newUpdate, middle, result)
	}

	return result, debug, nil
}

func pt2FixOrder(update []string, rules map[string][]string) []string {
	update = slices.Clone(update)
	slices.SortStableFunc(update, func(a, b string) int {
		rule, ok := rules[a]
		if !ok {
			return 0
		}
		if slices.Contains(rule, b) {
			return 1
		}
		return -1
	})
	return update
}
