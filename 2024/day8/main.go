package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	_ "embed"
)

const (
	VisualizeStep = "==========STEP==========\n"
	VisualizeData = "==========DATA==========\n"
)

// 991 too low

func main() {
	defaultInput := "input.txt"
	defaultPart := ""

	// defaultInput = "input-example-5.txt"
	// defaultPart = "2"

	var verboseDebug, outputFile bool
	var inputPath, partFilter string
	flag.StringVar(&inputPath, "input", defaultInput, "")
	flag.StringVar(&partFilter, "part", defaultPart, "")
	flag.BoolVar(&verboseDebug, "v", false, "verbose debug")
	flag.BoolVar(&outputFile, "o", false, "write output file")
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

		if outputFile {
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
		}

		debug.Close()
		slog.Info("finished running part", "func", funcName, "duration", duration, "result", result)
	}
}

func Part1(input string, debug *Debugger) (any, error) {
	result := 0

	lines := strings.Split(input, "\n")
	height := len(lines)
	width := len(lines[0])

	antenna := map[rune][]Pos{}
	antinodes := map[Pos]struct{}{}

	addAntinode := func(pos Pos) {
		if pos.Row >= 0 && pos.Row < height && pos.Col >= 0 && pos.Col < width {
			if _, ok := antinodes[pos]; !ok {
				antinodes[pos] = struct{}{}
				result++
				debug.WriteFormat("In bounds: %v\n", pos)
				return
			}

			debug.WriteFormat("Already in bounds: %v\n", pos)
			return
		}

		debug.WriteFormat("Out of bounds: %v\n", pos)
	}

	for row := range lines {
		for col, r := range []rune(lines[row]) {
			if r == '.' || r == '#' {
				continue
			}

			debug.WriteString(VisualizeStep)

			if locations, ok := antenna[r]; ok {
				for _, pos := range locations {
					diffRow := Abs(row - pos.Row)
					diffCol := Abs(col - pos.Col)
					mod := 1
					if col < pos.Col {
						mod = -1
					}

					addAntinode(Pos{Row: pos.Row - diffRow, Col: pos.Col - mod*diffCol})
					addAntinode(Pos{Row: row + diffRow, Col: col + mod*diffCol})
				}
			}
			antenna[r] = append(antenna[r], Pos{Row: row, Col: col})

			debug.WriteFormat("Result: %d\n%s", len(antinodes), VisualizeData)
			debug.WriteFunc(func() string { return fieldToString(height, width, antinodes, antenna) })
		}
	}

	return result, nil
}

func Part2(input string, debug *Debugger) (any, error) {
	result := 0

	lines := strings.Split(input, "\n")
	height := len(lines)
	width := len(lines[0])

	antenna := map[rune][]Pos{}
	antinodes := map[Pos]struct{}{}

	addAntinode := func(pos Pos) bool {
		if pos.Row >= 0 && pos.Row < height && pos.Col >= 0 && pos.Col < width {
			if _, ok := antinodes[pos]; !ok {
				antinodes[pos] = struct{}{}
				result++
				debug.WriteFormat("In bounds: %v\n", pos)
				return true
			}

			debug.WriteFormat("Already in bounds: %v\n", pos)
			return true
		}

		debug.WriteFormat("Out of bounds: %v\n", pos)
		return false
	}

	for row := range lines {
		for col, r := range []rune(lines[row]) {
			if r == '.' || r == '#' {
				continue
			}

			debug.WriteString(VisualizeStep)

			if locations, ok := antenna[r]; ok {
				for _, pos := range locations {
					diffRow := Abs(row - pos.Row)
					diffCol := Abs(col - pos.Col)
					mod := 1
					if col < pos.Col {
						mod = -1
					}

					count := 0
					up, down := true, true
					for up || down {
						if up {
							up = addAntinode(Pos{Row: pos.Row - count*diffRow, Col: pos.Col - count*mod*diffCol})
						}
						if down {
							down = addAntinode(Pos{Row: row + count*diffRow, Col: col + count*mod*diffCol})
						}
						count++
					}
				}
			}
			antenna[r] = append(antenna[r], Pos{Row: row, Col: col})

			debug.WriteFormat("Result: %d\n%s", len(antinodes), VisualizeData)
			debug.WriteFunc(func() string { return fieldToString(height, width, antinodes, antenna) })
			debug.Flush()
		}
	}

	return result, nil
}

type Pos struct {
	Row, Col int
}

func fieldToString(height, width int, antinodes map[Pos]struct{}, antenna map[rune][]Pos) string {
	var field [][]rune

	for row := 0; row < height; row++ {
		field = append(field, []rune(strings.Repeat(".", width)))
	}

	for pos := range antinodes {
		field[pos.Row][pos.Col] = '#'
	}

	for r, locations := range antenna {
		for _, pos := range locations {
			field[pos.Row][pos.Col] = r
		}
	}

	output := ""
	for _, row := range field {
		output += string(row) + "\n"
	}

	return output
}
