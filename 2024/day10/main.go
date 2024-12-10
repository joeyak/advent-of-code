package main

import (
	"flag"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"time"
)

var width, height int

func main() {
	defaultInput := "input.txt"
	defaultPart := ""

	// defaultInput = "input-example-1.txt"
	// defaultPart = "1"

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

var (
	Up    = NewPos(-1, 0)
	Down  = NewPos(1, 0)
	Left  = NewPos(0, -1)
	Right = NewPos(0, 1)
)

type Pos int

func NewPos(row, col int) Pos {
	return Pos(row*100 + col)
}

func (p Pos) Coords() (row int, col int) {
	return int(p) / 100, int(p) % 100
}

func Part1(input string, debug *Debugger) (any, error) {
	result := 0

	var trailheads []Pos
	field := map[Pos]int{}
	for row, line := range strings.Split(input, "\n") {
		height++
		if width == 0 {
			width = len(line)
		}

		for col, r := range line {
			if n := int(r - 48); n >= 0 && n < 10 {
				pos := NewPos(row, col)
				if n == 0 {
					trailheads = append(trailheads, pos)
					continue
				}
				field[pos] = n
			}
		}
	}

	for _, pos := range trailheads {
		result += len(followTrailPt1(debug, field, 0, pos, nil))

		debug.Flush()
	}

	return result, nil
}

func followTrailPt1(debug *Debugger, field map[Pos]int, val int, pos Pos, path []Pos) Hash[Pos] {
	ends := make(Hash[Pos])

	path = append(path, pos)
	debug.WriteFunc(debugPath(field, path))

	if val == 9 {
		ends.Add(pos)
		return ends
	}

	for _, mod := range []Pos{Up, Down, Left, Right} {
		newPos := pos + mod
		if newPos < 0 {
			continue
		}

		if newVal, ok := field[newPos]; ok && newVal == val+1 {
			newEnds := followTrailPt1(debug, field, newVal, newPos, path)
			maps.Copy(ends, newEnds)
		}
	}

	return ends
}

func Part2(input string, debug *Debugger) (any, error) {
	result := 0

	var trailheads []Pos
	field := map[Pos]int{}
	for row, line := range strings.Split(input, "\n") {
		height++
		if width == 0 {
			width = len(line)
		}

		for col, r := range line {
			if n := int(r - 48); n >= 0 && n < 10 {
				pos := NewPos(row, col)
				if n == 0 {
					trailheads = append(trailheads, pos)
					continue
				}
				field[pos] = n
			}
		}
	}

	for _, pos := range trailheads {
		result += followTrailPt2(debug, field, 0, pos, nil)

		debug.Flush()
	}

	return result, nil
}

func followTrailPt2(debug *Debugger, field map[Pos]int, val int, pos Pos, path []Pos) int {
	count := 0

	ends := make(Hash[Pos])

	path = append(path, pos)
	debug.WriteFunc(debugPath(field, path))

	if val == 9 {
		ends.Add(pos)
		return 1
	}

	for _, mod := range []Pos{Up, Down, Left, Right} {
		newPos := pos + mod
		if newPos < 0 {
			continue
		}

		if newVal, ok := field[newPos]; ok && newVal == val+1 {
			count += followTrailPt2(debug, field, newVal, newPos, path)
		}
	}

	return count
}

func debugPath(field map[Pos]int, path []Pos) func() string {
	return func() string {
		meta := slices.Repeat([]rune(strings.Repeat(".", width)+"\n"), height)
		data := slices.Repeat([]rune(strings.Repeat(".", width)+"\n"), height)
		for _, pos := range path {
			row, col := pos.Coords()
			n := (row)*(width+1) + col
			data[n] = '#'
			meta[n] = rune(field[pos] + 48)
		}

		return VisualizeStep + string(meta) + VisualizeData + string(data)
	}
}
