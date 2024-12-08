package main

import (
	"flag"
	"fmt"
	"log/slog"
	"math"
	"math/big"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	_ "embed"
)

func main() {
	var verboseDebug bool
	var inputPath, partFilter string
	flag.StringVar(&inputPath, "input", "input.txt", "")
	flag.StringVar(&partFilter, "part", "", "")
	flag.BoolVar(&verboseDebug, "v", false, "verbose debug")
	flag.Parse()

	// inputPath = "input-example-1.txt"
	// partFilter = "2"
	// verboseDebug = true

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

		debug.Close()

		slog.Info("finished running part", "func", funcName, "duration", duration, "result", result)
	}
}

func Part1(input string, debug *Debugger) (any, error) {
	result := 0

	for ln, line := range strings.Split(input, "\n") {
		parts := strings.Split(line, ":")
		expected := parts[0]
		var nums []int
		for _, s := range strings.Split(strings.TrimSpace(parts[1]), " ") {
			num, _ := strconv.Atoi(s)
			nums = append(nums, num)
		}

		added := false
		debug.WriteString(line + "\n")
		pow := len(nums) - 1
		opLen := 2
		loopMax := int(math.Pow(float64(opLen), float64(pow)))

		fmt.Printf("\r%d max(%d) opLen(%d) pow(%d)", ln+1, loopMax, opLen, pow)
		for i := 0; i < loopMax; i++ {
			ops := fmt.Sprintf(fmt.Sprintf("%%0%ds", pow), big.NewInt(int64(i)).Text(opLen))
			debug.WriteString(ops + " | ")

			actual := recurseApplyOperators([]rune(ops), nums, "", debug)

			debug.WriteString(fmt.Sprintf("= %s [%s]", actual, expected))

			if actual == expected && !added {
				n, _ := strconv.Atoi(expected)
				result += n
				added = true
				debug.WriteString(fmt.Sprintf(" => %t", actual == expected))
			}
			debug.WriteString("\n")
		}
		debug.WriteString("\n")
	}
	fmt.Print("\r")

	return result, nil
}

func Part2(input string, debug *Debugger) (any, error) {
	result := 0

	for ln, line := range strings.Split(input, "\n") {
		parts := strings.Split(line, ":")
		expected := parts[0]
		var nums []int
		for _, s := range strings.Split(strings.TrimSpace(parts[1]), " ") {
			num, _ := strconv.Atoi(s)
			nums = append(nums, num)
		}

		added := false
		debug.WriteString(line + "\n")
		pow := len(nums) - 1
		opLen := 3
		loopMax := int(math.Pow(float64(opLen), float64(pow)))

		fmt.Printf("\r%d max(%d) opLen(%d) pow(%d)", ln+1, loopMax, opLen, pow)
		for i := 0; i < loopMax; i++ {
			// Added break cause almost ran out of memory concatenating strings
			if added {
				break
			}

			ops := fmt.Sprintf(fmt.Sprintf("%%0%ds", pow), big.NewInt(int64(i)).Text(opLen))
			debug.WriteString(ops + " | ")

			actual := recurseApplyOperators([]rune(ops), nums, "", debug)

			debug.WriteString(fmt.Sprintf("= %s [%s]", actual, expected))

			if actual == expected && !added {
				n, _ := strconv.Atoi(expected)
				result += n
				added = true
				debug.WriteString(fmt.Sprintf(" => %t", actual == expected))
			}
			debug.WriteString("\n")
		}
		debug.WriteString("\n")
	}
	fmt.Print("\r")

	return result, nil
}

func recurseApplyOperators(ops []rune, data []int, actual string, debug *Debugger) string {
	if len(ops) == 0 {
		return actual
	}

	if actual == "" {
		actual = strconv.Itoa(data[0])
		debug.WriteString(actual + " ")
		return recurseApplyOperators(ops, data[1:], actual, debug)
	}

	numStr := strconv.Itoa(data[0])
	actualNum, _ := strconv.Atoi(actual)
	switch ops[0] {
	case '0':
		debug.WriteString("+ ")
		actual = strconv.Itoa(actualNum + data[0])
	case '1':
		debug.WriteString("* ")
		actual = strconv.Itoa(actualNum * data[0])
	case '2':
		debug.WriteString("|| ")
		actual += numStr
	}
	debug.WriteString(fmt.Sprintf("%s (%s) ", numStr, actual))

	return recurseApplyOperators(ops[1:], data[1:], actual, debug)
}
