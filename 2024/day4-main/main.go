package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"runtime"
	"slices"
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

type Direction int

const (
	Up Direction = (1 << iota)
	Down
	Left
	Right
)

func Part1(input string) (any, string, error) {
	result := 0
	debug := ""

	var runes [][]rune
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		runes = append(runes, []rune(line))
	}

	var goodCoords [][]int
	var debugRunes [][]rune
	for y := 0; y < len(runes); y++ {
		debugRunes = append(debugRunes, []rune(strings.Repeat(".", len(runes[y]))))
		for x := 0; x < len(runes[y]); x++ {
			if runes[y][x] == 'X' {
				for _, dir := range []Direction{Up, Down, Left, Right, Up | Left, Up | Right, Down | Left, Down | Right} {
					if coords := pt1CheckDirection(runes, y, x, 0, dir, nil); coords != nil {
						goodCoords = append(goodCoords, coords...)
						result++
					}
				}
			}
		}
	}

	for _, coords := range goodCoords {
		debugRunes[coords[0]][coords[1]] = runes[coords[0]][coords[1]]
	}

	for _, runes := range debugRunes {
		debug += string(runes) + "\n"
	}

	return result, debug, nil
}

func pt1CheckDirection(runes [][]rune, y, x, step int, dir Direction, goodCoords [][]int) [][]int {
	if c := rune("XMAS"[step]); c == runes[y][x] {
		goodCoords = append(slices.Clone(goodCoords), []int{y, x})
		if c == 'S' {
			return goodCoords
		}

		if dir&Up != 0 {
			y--
		}
		if dir&Down != 0 {
			y++
		}
		if dir&Left != 0 {
			x--
		}
		if dir&Right != 0 {
			x++
		}

		if y >= 0 && y < len(runes) && x >= 0 && x < len(runes[0]) {
			return pt1CheckDirection(runes, y, x, step+1, dir, goodCoords)
		}
	}
	return nil
}

func Part2(input string) (any, string, error) {
	result := 0
	debug := ""

	var runes [][]rune
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		runes = append(runes, []rune(line))
	}

	var goodCoords [][]int
	var debugRunes [][]rune

	// Add two extra for skipped row checks at start and end
	debugRunes = append(debugRunes, []rune(strings.Repeat(".", len(runes[0]))))
	debugRunes = append(debugRunes, []rune(strings.Repeat(".", len(runes[0]))))

	for y := 1; y < len(runes)-1; y++ {
		debugRunes = append(debugRunes, []rune(strings.Repeat(".", len(runes[y]))))
		for x := 1; x < len(runes[y])-1; x++ {
			if runes[y][x] == 'A' {
				upLeft := runes[y-1][x-1]
				upRight := runes[y-1][x+1]
				downLeft := runes[y+1][x-1]
				downRight := runes[y+1][x+1]

				if ((upLeft == 'M' && downRight == 'S') || (upLeft == 'S' && downRight == 'M')) &&
					((upRight == 'M' && downLeft == 'S') || (upRight == 'S' && downLeft == 'M')) {
					result++
					goodCoords = append(goodCoords, []int{y, x})
				}
			}
		}
	}

	for _, coords := range goodCoords {
		for _, mod := range [][]int{
			{-1, -1},
			{-1, +1},
			{+1, -1},
			{+1, +1},
			{0, 0},
		} {
			y := coords[0] + mod[0]
			x := coords[1] + mod[1]
			debugRunes[y][x] = runes[y][x]
		}
	}

	for _, runes := range debugRunes {
		debug += string(runes) + "\n"
	}

	return result, debug, nil
}
