package main

import (
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

//go:embed input.txt
var input string

func main() {
	for _, f := range []func() (any, string, error){Part1, Part2} {
		funcName := strings.Split(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), ".")[1]

		start := time.Now()
		result, debug, err := f()
		duration := time.Since(start)

		if err != nil {
			slog.Error("could not run part", "func", funcName, "err", err)
			break
		}

		outFile, err := os.OpenFile(fmt.Sprintf("output-%s.txt", funcName), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
		if err != nil {
			slog.Error("could not create or trunc output file", "err", err)
			os.Exit(1)
		}
		defer outFile.Close()

		_, err = outFile.WriteString(fmt.Sprintf("%v", result))
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

func Part1() (any, string, error) {
	debug := ""
	result := 0
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		var dataLine []int
		var diffs []int

		safe := true
		increasing := false
		decreasing := false
		lastNum := 0
		for i, numString := range strings.Split(line, " ") {
			num, _ := strconv.Atoi(numString)
			dataLine = append(dataLine, num)
			if i != 0 {
				diff := lastNum - num
				diffs = append(diffs, diff)

				if diff > 0 {
					increasing = true
				}

				if diff < 0 {
					decreasing = true
				}

				if (diff < -3 || diff == 0 || diff > 3) || (diff > 0 && decreasing) || (diff < 0 && increasing) {
					safe = false
					break
				}
			}
			lastNum = num
		}

		if safe {
			result++
		}
		debug += fmt.Sprintf("data:%v diffs:%v safe(%t)\n", dataLine, diffs, safe)
	}

	return result, debug, nil
}

func Part2() (any, string, error) {
	debug := ""
	result := 0

	var data [][]int
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		var dataLine []int
		for _, numString := range strings.Split(line, " ") {
			num, _ := strconv.Atoi(numString)
			dataLine = append(dataLine, num)
		}
		data = append(data, dataLine)
	}

	for _, dataLine := range data {
		failed, usedData, diffs := comparePt2(dataLine, true)
		if failed == -1 {
			result++
		}
		debug += fmt.Sprintf("orig:%v used:%v diffs:%v safe(%t", dataLine, usedData, diffs, failed == -1)
		if failed != -1 {
			debug += "|" + strconv.Itoa(failed+1)
		}
		debug += ")\n"
	}

	return result, debug, nil
}

func comparePt2(data []int, first bool) (int, []int, any) {
	failedColumn := -1

	lastDiff := 0
	var diffs []int
	for i := range data {
		if i == 0 {
			continue
		}
		diff := data[i] - data[i-1]
		diffs = append(diffs, diff)

		if (i == 1 && diff == 0) || (i > 1 && diff*lastDiff <= 0) || diff > 3 || diff < -3 {
			failedColumn = i
			break
		}

		lastDiff = diff
	}

	if failedColumn != -1 {
		if first {
			failed, newData, diffs2 := comparePt2(append(slices.Clone(data[:failedColumn]), data[failedColumn+1:]...), false)
			if failed != -1 {
				failed, newData, diffs2 = comparePt2(append(slices.Clone(data[:failedColumn-1]), data[failedColumn:]...), false)
				if failed != -1 {
					failed, newData, diffs2 = comparePt2(data[1:], false)
				}
			}
			return failed, newData, [][]int{diffs, diffs2.([]int)}
		}
		return failedColumn, data, diffs
	}
	return failedColumn, data, diffs
}
