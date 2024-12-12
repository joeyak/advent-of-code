package main

import (
	"flag"
	"fmt"
	"iter"
	"log/slog"
	"os"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"time"
)

func main() {
	defaultInput := "input.txt"
	defaultPart := ""

	// defaultInput = "input-example-1.txt"
	// defaultPart = "1"

	var verboseDebug bool
	var inputPath, partFilter string
	flag.StringVar(&inputPath, "input", defaultInput, "")
	flag.StringVar(&partFilter, "part", defaultPart, "")
	flag.BoolVar(&verboseDebug, "v", false, "verbose debug")
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

type Crop struct {
	Region int
	Sides  int
	Rune   rune
}

func Part1(input string, debug *Debugger) (any, error) {
	result := 0

	lines := strings.Split(input, "\n")
	height := len(lines)
	width := len(lines[0])
	field := make([][]*Crop, height)
	for row, line := range lines {
		field[row] = make([]*Crop, width)
		for col, r := range line {
			field[row][col] = &Crop{Rune: r}
		}
	}

	debug.WriteFunc(func() string {
		return "Orig" + strings.Repeat(" ", max(1, width-2)) + "Sides" + strings.Repeat(" ", max(1, width-3)) + "Region\n"
	})

	regions := map[int][]*Crop{}
	regionCounter := 0
	for row := range field {
		for dirIdx, compares := range [][]any{{slices.Backward[[]*Crop], slices.Backward[[][]int]}, {slices.All[[]*Crop], slices.All[[][]int]}} {
			dirField := compares[0].(func([]*Crop) iter.Seq2[int, *Crop])
			dirInt := compares[1].(func([][]int) iter.Seq2[int, []int])
			for col, crop := range dirField(field[row]) {
				for _, pos := range dirInt([][]int{{-1, 0}, {0, -1}, {0, +1}, {+1, 0}}) {
					checkRow, checkCol := row+pos[0], col+pos[1]
					if checkRow < 0 || checkRow >= height || checkCol < 0 || checkCol >= width {
						if dirIdx == 1 {
							crop.Sides++
						}
						continue
					}

					check := field[checkRow][checkCol]
					if crop.Region == 0 && check.Region != 0 && crop.Rune == check.Rune {
						crop.Region = check.Region
						continue
					}

					if dirIdx == 1 {
						if crop.Rune != check.Rune {
							crop.Sides++
						}
					}
				}

				if dirIdx == 1 {
					if crop.Region == 0 {
						regionCounter++
						crop.Region = regionCounter
					}

					regions[crop.Region] = append(regions[crop.Region], crop)
				}
			}
		}

		debug.WriteFunc(func() string {
			s := append([]rune(strings.Repeat(" ", width*3+4)), '\n')

			for col, crop := range field[row] {
				s[col] = crop.Rune
				s[col+width+2] = rune(crop.Sides + 48)
				s[col+width*2+4] = rune(crop.Region + 48)
			}

			return string(s)
		})
		debug.Flush()
	}

	debug.WriteString("\n")

	for region := 1; region <= regionCounter; region++ {
		crops := regions[region]

		area := len(crops)
		perimeter := 0
		for i := range crops {
			perimeter += crops[i].Sides
		}
		result += area * perimeter

		debug.WriteFormat("[%d|%s] a(%d) + p(%d) = c(%d) => %d\n", region, string(crops[0].Rune), area, perimeter, area*perimeter, result)
	}

	return result, nil
}

func Part2(input string, debug *Debugger) (any, error) {
	result := 0

	for _, line := range strings.Split(input, "\n") {
		_ = line
	}

	return result, nil
}
