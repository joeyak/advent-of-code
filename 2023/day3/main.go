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
	"unicode"

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

	var symbols []Pos
	var numbers []Number
	for y, line := range iterLines(input) {
		number := Number{}
		for x, r := range line {
			if unicode.IsDigit(r) {
				if number.Value == "" {
					number = Number{Pos: Pos{X: x, Y: y}}
				}
				number.Value += string(r)
				continue
			}

			if number.Value != "" {
				numbers = append(numbers, number)
			}

			if r != '.' {
				symbols = append(symbols, Pos{X: x, Y: y})
			}
			number = Number{}
		}
		if number.Value != "" {
			numbers = append(numbers, number)
		}
	}

	for _, n := range numbers {
		isPartNum := false
		for _, s := range symbols {
			for mod := range n.Value {
				diffX := s.X - (n.X + mod)
				diffY := s.Y - n.Y
				if diffX <= 1 && diffX >= -1 && diffY <= 1 && diffY >= -1 {
					isPartNum = true
					break
				}
			}
			if isPartNum {
				break
			}
		}

		if isPartNum {
			num, _ := strconv.Atoi(n.Value)
			result += num
		}

		debug += fmt.Sprintf("Number: %+v - %t\n", n, isPartNum)
	}

	return result, debug, nil
}

func Part2(input string) (any, string, error) {
	result := 0
	debug := ""

	var gears []*Gear
	var parts []*Part
	for y, line := range iterLines(input) {
		part := &Part{}
		for x, r := range line {
			if unicode.IsDigit(r) {
				if part.Value == "" {
					part = &Part{Pos: Pos{X: x, Y: y}}
				}
				part.Value += string(r)
				continue
			}

			if part.Value != "" {
				parts = append(parts, part)
			}

			if r != '.' {
				gears = append(gears, &Gear{Pos: Pos{X: x, Y: y}, Value: r})
			}
			part = &Part{}
		}
		if part.Value != "" {
			parts = append(parts, part)
		}
	}

	for _, part := range parts {
		for _, gear := range gears {
			for mod := range part.Value {
				diffX := gear.X - (part.X + mod)
				diffY := gear.Y - part.Y
				if diffX <= 1 && diffX >= -1 && diffY <= 1 && diffY >= -1 {
					part.Gears = append(part.Gears, gear)
					gear.Parts = append(gear.Parts, part)
					break
				}
			}
		}
		debug += fmt.Sprintf("Part: %+v\n", part)
	}

	for _, gear := range gears {
		debug += fmt.Sprintf("Gear: %v %v", gear.Pos, gear.Value)
		if gear.Value == '*' {
			if len(gear.Parts) == 2 {
				debug += " -"
				gearValue := 1
				for _, part := range gear.Parts {
					debug += fmt.Sprintf(" %v", part.Value)
					num, _ := strconv.Atoi(part.Value)
					gearValue *= num
				}
				result += gearValue
				debug += fmt.Sprintf(" - %d", gearValue)
			}
		}
		debug += "\n"
	}

	return result, debug, nil
}

type Pos struct {
	X, Y int
}

type Number struct {
	Pos
	Value string
}

type Part struct {
	Pos
	Value string
	Gears []*Gear
}

type Gear struct {
	Pos
	Value rune
	Parts []*Part
}
