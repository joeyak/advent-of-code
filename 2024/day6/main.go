package main

import (
	"encoding/json"
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

const VisualizeStep = "==========STEP==========\n"

var (
	verboseDebug    bool
	DirectionChange = map[string]string{
		"^": ">",
		">": "V",
		"V": "<",
		"<": "^",
	}
	DirectionSymbol = map[string]string{
		"^": "|",
		"V": "|",
		">": "-",
		"<": "-",
	}
)

func main() {
	var inputPath, partFilter string
	flag.StringVar(&inputPath, "input", "input.txt", "")
	flag.StringVar(&partFilter, "part", "", "")
	flag.BoolVar(&verboseDebug, "v", false, "verbose debug")
	flag.Parse()

	// inputPath = "input-example-1.txt"

	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		slog.Error("could not read file", "path", inputPath, "err", err)
		os.Exit(1)
	}

	input := strings.TrimSuffix(strings.ReplaceAll(string(inputData), "\r\n", "\n"), "\n")

	for _, f := range []func(string) (any, string, error){Part1, Part2} {
		funcName := strings.Split(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), ".")[1]
		if !strings.HasSuffix(funcName, partFilter) {
			continue
		}

		start := time.Now()
		result, debug, err := f(input)
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

func Part1(input string) (any, string, error) {
	result := 0
	debug := ""

	var state State
	for i, r := range input {
		x := i - (state.Height * (state.Width + 1))

		switch r {
		case '\n':
			if state.Width == 0 {
				state.Width = i
			}
			state.Height++
		case '^', 'V', '<', '>':
			state.Guard = Guard{
				Pos: Pos{X: x, Y: state.Height},
				Dir: string(r),
			}
		case '#':
			state.Obstacles = append(state.Obstacles, Pos{X: x, Y: state.Height})
		}
	}
	state.Height++

	b, _ := json.MarshalIndent(state, "", "    ")
	debug += string(b) + "\n"

	maxStep := -1
	for state.Step < maxStep || maxStep == -1 {
		state.Step++
		fmt.Printf("\rStep %d", state.Step)
		if state.GuardOutOfBounds() {
			fmt.Print("\r")
			slog.Info("out of bounds", "func", "Part1", "step", state.Step, "X", state.Guard.Pos.X, "Y", state.Guard.Pos.Y)
			break
		}
		if verboseDebug {
			debug += VisualizeStep + state.Debug()
		}

		newPos := state.Guard.Pos
		switch state.Guard.Dir {
		case "^":
			newPos.Y--
		case "V":
			newPos.Y++
		case "<":
			newPos.X--
		case ">":
			newPos.X++
		}

		hit := false
		for _, obstacle := range state.Obstacles {
			if newPos.X == obstacle.X && newPos.Y == obstacle.Y {
				hit = true
				break
			}
		}

		state.Path = append(state.Path, state.Guard)
		if hit {
			state.Guard.Dir = DirectionChange[state.Guard.Dir]
			continue
		}
		state.Guard.Pos = newPos
	}
	lastMapState := state.Debug()
	debug += VisualizeStep + lastMapState

	for _, r := range lastMapState {
		if r != '#' && r != '.' && r != ' ' && r != '\n' {
			result++
		}
	}

	return result, debug, nil
}

func Part2(input string) (any, string, error) {
	result := 0
	debug := ""

	return result, debug, nil
}

type Pos struct {
	X, Y int
}

type Guard struct {
	Pos
	Dir string
}

type State struct {
	Step      int
	Width     int
	Height    int
	Guard     Guard
	Path      []Guard
	Obstacles []Pos
}

func (s State) GuardOutOfBounds() bool {
	return s.Guard.X < 0 || s.Guard.X >= s.Width || s.Guard.Y < 0 || s.Guard.Y >= s.Height
}

func (s State) Debug() string {
	var builder strings.Builder
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			if s.Guard.X == x && s.Guard.Y == y {
				builder.WriteString(s.Guard.Dir)
				continue
			}

			hasObstacle := false
			for _, obstacle := range s.Obstacles {
				if obstacle.X == x && obstacle.Y == y {
					hasObstacle = true
					break
				}
			}
			if hasObstacle {
				builder.WriteString("#")
				continue
			}

			pathSymbol := ""
			for i, guard := range s.Path {
				if guard.X == x && guard.Y == y {
					if i == 0 {
						pathSymbol = "@"
						break
					}

					symbol := DirectionSymbol[guard.Dir]
					if pathSymbol != "" && pathSymbol != symbol {
						pathSymbol = "+"
						continue
					}

					pathSymbol = symbol
				}
			}
			if pathSymbol != "" {
				builder.WriteString(pathSymbol)
				continue
			}

			builder.WriteString(".")
		}
		builder.WriteString("\n")
	}
	return builder.String()
}
