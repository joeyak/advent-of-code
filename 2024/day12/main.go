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
)

func main() {
	defaultInput := "input.txt"
	defaultPart := ""

	// defaultInput = "input-example-3.txt"
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
			field[row][col] = &Crop{Rune: r, Region: -1}
		}
	}

	regionCounter := 0
	regions := map[int][]*Crop{}
	for row := range field {
		for col := range field[row] {
			if field[row][col].Region == -1 {
				regions[regionCounter] = pathCrops(field, Hash[*Crop]{}, row, col, width, height)
				for _, crop := range regions[regionCounter] {
					crop.Region = regionCounter
				}

				debug.WriteFunc(func() string {
					s := VisualizeStep + "Orig" + strings.Repeat(" ", max(1, width-2)) + "Region" + strings.Repeat(" ", max(1, width-4)) + "Sides\n"
					for _, row := range field {
						ss := []rune(strings.Repeat(" ", width*3+4))

						for col, crop := range row {
							if crop.Region == -1 {
								ss[col] = '.'
								ss[col+width+2] = '.'
								ss[col+width*2+4] = '.'
							} else {
								ss[col] = crop.Rune
								ss[col+width+2] = rune(crop.Region + 48)
								ss[col+width*2+4] = rune(crop.Sides + 48)
							}
						}

						s += string(ss) + "\n"
					}
					return s
				})
				debug.Flush()
				regionCounter++
			}
		}
	}

	debug.WriteString(VisualizeEnd)

	for region := 0; region < len(regions); region++ {
		crops := regions[region]

		area := len(crops)
		perimeter := 0
		var r rune
		for _, crop := range crops {
			perimeter += crop.Sides
			r = crop.Rune
		}
		result += area * perimeter

		debug.WriteFormat("[%d|%s] a(%d) + p(%d) = c(%d) => %d\n", region, string(r), area, perimeter, area*perimeter, result)
	}

	return result, nil
}

func pathCrops(field [][]*Crop, checked Hash[*Crop], row, col, width, height int) []*Crop {
	crop := field[row][col]

	checked.Add(crop)
	crops := []*Crop{crop}

	for _, dir := range [][]int{{-1, 0}, {0, -1}, {1, 0}, {0, 1}} {
		cRow, cCol := row+dir[0], col+dir[1]
		if cRow < 0 || cRow >= height || cCol < 0 || cCol >= width {
			crop.Sides++
			continue
		}

		checkCrop := field[cRow][cCol]
		if checked.Has(checkCrop) {
			continue
		}

		if crop.Rune != checkCrop.Rune {
			crop.Sides++
			continue
		}

		crops = append(crops, pathCrops(field, checked, cRow, cCol, width, height)...)
	}

	return crops
}

func Part2(input string, debug *Debugger) (any, error) {
	result := 0

	for _, line := range strings.Split(input, "\n") {
		_ = line
	}

	return result, nil
}
