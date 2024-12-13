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

var veryVerbose bool

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
	flag.BoolVar(&veryVerbose, "vv", false, "very verbose debug")
	flag.Parse()

	if veryVerbose {
		verboseDebug = true
	}

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

const (
	Up Side = 1 << iota
	Down
	Left
	Right
)

type Side byte

func (s Side) Count() int {
	c := 0

	for _, dir := range []Side{Up, Down, Left, Right} {
		if s&dir != 0 {
			c++
		}
	}

	return c
}

type Crop struct {
	Row, Col int
	Region   int
	Sides    Side
	Rune     rune
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
			field[row][col] = &Crop{
				Row:    row,
				Col:    col,
				Rune:   r,
				Region: -1,
			}
		}
	}

	fireDebug := func() {
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
						ss[col+width*2+4] = rune(crop.Sides.Count() + 48)
					}
				}

				s += string(ss) + "\n"
			}
			return s
		})
		debug.Flush()
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

				if veryVerbose {
					fireDebug()
				}

				regionCounter++
			}
		}
	}

	if !veryVerbose {
		fireDebug()
	}

	debug.WriteString(VisualizeEnd)

	for region := 0; region < len(regions); region++ {
		crops := regions[region]

		area := len(crops)
		perimeter := 0
		for _, crop := range crops {
			perimeter += crop.Sides.Count()
		}
		result += area * perimeter

		debug.WriteFormat("[%d|%s] a(%d) + p(%d) = c(%d) => %d\n", region, string(crops[0].Rune), area, perimeter, area*perimeter, result)
	}

	return result, nil
}

type Direction struct {
	Side     Side
	Row, Col int
}

func pathCrops(field [][]*Crop, checked Hash[*Crop], row, col, width, height int) []*Crop {
	crop := field[row][col]

	checked.Add(crop)
	crops := []*Crop{crop}

	for _, dir := range []Direction{
		{Side: Up, Row: -1, Col: 0},
		{Side: Left, Row: 0, Col: -1},
		{Side: Down, Row: 1, Col: 0},
		{Side: Right, Row: 0, Col: 1},
	} {
		cRow, cCol := row+dir.Row, col+dir.Col
		if cRow < 0 || cRow >= height || cCol < 0 || cCol >= width {
			crop.Sides |= dir.Side
			continue
		}

		checkCrop := field[cRow][cCol]
		if checked.Has(checkCrop) {
			continue
		}

		if crop.Rune != checkCrop.Rune {
			crop.Sides |= dir.Side
			continue
		}

		crops = append(crops, pathCrops(field, checked, cRow, cCol, width, height)...)
	}

	return crops
}

func Part2(input string, debug *Debugger) (any, error) {
	result := 0

	lines := strings.Split(input, "\n")
	height := len(lines)
	width := len(lines[0])

	field := make([][]*Crop, height)
	for row, line := range lines {
		field[row] = make([]*Crop, width)
		for col, r := range line {
			field[row][col] = &Crop{
				Row:    row,
				Col:    col,
				Rune:   r,
				Region: -1,
			}
		}
	}

	fireDebug := func() {
		debug.WriteFunc(func() string {
			s := VisualizeStep + "Orig" + strings.Repeat(" ", max(1, width-2)) + "Region\n"
			for _, row := range field {
				ss := []rune(strings.Repeat(" ", width*2+2))

				for col, crop := range row {
					if crop.Region == -1 {
						ss[col] = '.'
						ss[col+width+2] = '.'
					} else {
						ss[col] = crop.Rune
						ss[col+width+2] = rune(crop.Region + 48)
					}
				}

				s += string(ss) + "\n"
			}
			return s
		})
		debug.Flush()
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

				if veryVerbose {
					fireDebug()
				}

				regionCounter++
			}
		}
	}

	if !veryVerbose {
		fireDebug()
	}

	debug.WriteString(VisualizeEnd)

	for region := 0; region < len(regions); region++ {
		crops := regions[region]

		area := len(crops)
		perimeter := bulkPerimeterFromRegion(crops, width, height)
		result += area * perimeter

		debug.WriteFormat("[%d|%s] a(%d) + p(%d) = c(%d) => %d\n", region, string(crops[0].Rune), area, perimeter, area*perimeter, result)
	}

	return result, nil
}

func bulkPerimeterFromRegion(cropList []*Crop, width, height int) int {
	_ = cropList
	_ = width
	_ = height
	return 0
	// keyMod := int(math.Pow(10, math.Log10(float64(width))+1))
	// cropKey := func(row, col int) int {
	// 	return row*keyMod+col
	// }

	// crops := map[int]*Crop{}
	// for _, crop := range cropList {
	// 	crops[cropKey(crop.Row, crop.Col)] = crop
	// }

	// sides := 0
	// inFence := false
	// checked := 0

	// // horizontal fences
	// var upSides, downSides int
	// for row := 0; row < height && checked < height; row++ {
	// 	for col := 0; col < width-1 && checked < width; col++ {
	// 		crop, ok := crops[cropKey(row, col)]
	// 		if !ok {
	// 			inFence =false
	// 			continue
	// 		}

	// 	}
	// }

	// return sides
}
