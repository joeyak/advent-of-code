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
)

var maxSteps, maxPt2Steps int

func main() {
	defaultInput := "input.txt"
	defaultPart := ""
	defaultMaxSteps := 25
	defaultPt2MaxSteps := 75

	// defaultInput = "input-example-2.txt"
	// defaultPt2MaxSteps = 6
	// defaultPart = "2"

	var verboseDebug bool
	var inputPath, partFilter string
	flag.StringVar(&inputPath, "input", defaultInput, "")
	flag.StringVar(&partFilter, "part", defaultPart, "")
	flag.BoolVar(&verboseDebug, "v", false, "verbose debug")
	flag.IntVar(&maxSteps, "steps", defaultMaxSteps, "max steps to iterate")
	flag.IntVar(&maxPt2Steps, "steps2", defaultPt2MaxSteps, "max steps to iterate")
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
	data := strings.Split(input, " ")
	debug.WriteFunc(debugData(data))

	for i := 0; i < maxSteps; i++ {
		fmt.Printf("\r%d - %d", i, len(data))
		var newData []string
		for _, item := range data {
			if item == "0" {
				newData = append(newData, "1")
				continue
			}

			if len(item)%2 == 0 {
				for _, s := range []string{item[:len(item)/2], item[len(item)/2:]} {
					for s[0] == '0' && len(s) > 1 {
						s = s[1:]
					}
					newData = append(newData, s)
				}
				continue
			}

			num, _ := strconv.Atoi(item)
			newData = append(newData, strconv.Itoa(num*2024))
		}

		data = newData
		debug.WriteFunc(debugData(data))
	}
	fmt.Print("\r")

	return len(data), nil
}

func debugData[T map[string]int | []string](data T) func() string {
	return func() string {
		mapData, ok := any(data).(map[string]int)
		if !ok {
			mapData = map[string]int{}
			for _, s := range any(data).([]string) {
				mapData[s]++
			}
		}

		return fmt.Sprintf("%v\n", mapData)
	}
}

func Part2(input string, debug *Debugger) (any, error) {
	result := 0

	data := map[string]int{}
	for _, s := range strings.Split(input, " ") {
		data[s]++
	}
	debug.WriteFunc(debugData(data))

	for i := 0; i < maxPt2Steps; i++ {
		fmt.Printf("\r%d", i)

		newData := map[string]int{}
		for item, count := range data {
			if item == "0" {
				newData["1"] += count
				continue
			}

			if len(item)%2 == 0 {
				for _, s := range []string{item[:len(item)/2], item[len(item)/2:]} {
					for s[0] == '0' && len(s) > 1 {
						s = s[1:]
					}
					newData[s] += count
				}
				continue
			}

			num, _ := strconv.Atoi(item)
			newData[strconv.Itoa(num*2024)] += count
		}

		data = newData
		debug.WriteFunc(debugData(data))
	}
	fmt.Print("\r")

	for _, n := range data {
		result += n
	}

	return result, nil
}
