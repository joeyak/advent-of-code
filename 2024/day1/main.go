package main

import (
	"fmt"
	"log/slog"
	"math"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	_ "embed"
	"slices"
)

//go:embed input.txt
var input string

func main() {
	for _, f := range []func() (any, error){Part1, Part2} {
		funcName := strings.Split(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), ".")[1]

		start := time.Now()
		result, err := f()
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

		slog.Info("finished running part", "func", funcName, "duration", duration, "result", result)
	}
}

func Part1() (any, error) {
	var list1, list2 []int
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}

		numbers := strings.Split(strings.TrimSpace(line), "   ")
		num1, _ := strconv.Atoi(numbers[0])
		list1 = append(list1, num1)

		num2, _ := strconv.Atoi(numbers[1])
		list2 = append(list2, num2)
	}

	slices.Sort(list1)
	slices.Sort(list2)

	result := 0
	for i, v := range list1 {
		result += int(math.Abs(float64(list2[i] - v)))
	}

	return result, nil
}

func Part2() (any, error) {
	var list1 []int
	list2 := map[int]int{}
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}

		numbers := strings.Split(strings.TrimSpace(line), "   ")
		num1, _ := strconv.Atoi(numbers[0])
		list1 = append(list1, num1)

		num2, _ := strconv.Atoi(numbers[1])
		list2[num2]++
	}

	result := 0
	for _, v1 := range list1 {
		v2, ok := list2[v1]
		if ok {
			result += v1 * v2
		}
	}

	return result, nil
}
