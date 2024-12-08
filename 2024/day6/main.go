package main

import (
	"encoding/json"
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

const VisualizeStep = "==========STEP==========\n"

var (
	veryVerboseDebug bool

	DirectionChange = map[rune]rune{
		'^': '>',
		'>': 'V',
		'V': '<',
		'<': '^',
	}
	DirectionSymbol = map[rune]string{
		'^': "|",
		'V': "|",
		'>': "-",
		'<': "-",
	}
	pt1Result = -1
)

func main() {
	var inputPath, partFilter string
	var verboseDebug bool
	flag.StringVar(&inputPath, "input", "input.txt", "")
	flag.StringVar(&partFilter, "part", "", "")
	flag.BoolVar(&verboseDebug, "v", false, "verbose debug")
	flag.BoolVar(&veryVerboseDebug, "vv", false, "very verbose debug")
	flag.Parse()

	if veryVerboseDebug {
		verboseDebug = true
	}

	// inputPath = "input-example-1.txt"
	// partFilter = "2"
	// verboseDebug = false

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

		debug.Flush()
		slog.Info("finished running part", "func", funcName, "duration", duration, "result", result)
	}
}

func Part1(input string, debug *Debugger) (any, error) {
	result := 0

	state := State{LoopCheck: map[int]int{}}
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
				Dir: r,
			}
		case '#':
			state.Obstacles = append(state.Obstacles, Obstacle{Pos: Pos{X: x, Y: state.Height}, Symbol: '#'})
		}
	}
	state.Height++

	b, _ := json.MarshalIndent(state, "", "    ")
	debug.WriteString(string(b) + "\n")

	for !state.GuardOutOfBounds() {
		fmt.Printf("\rStep: %d", state.StepCount)
		state = state.Step()

		if (state.Width == 10 || (state.Width > 10 && state.StepCount%25 == 0)) || veryVerboseDebug {
			debug.WriteFunc(func() string { return VisualizeStep + state.Debug() })
		}
	}
	fmt.Print("\r")
	slog.Info("out of bounds", "func", "Part1", "step", state.StepCount, "X", state.Guard.Pos.X, "Y", state.Guard.Pos.Y, "W", state.Width, "H", state.Height)

	lastMapState := state.Debug()
	debug.WriteString(VisualizeStep + lastMapState)

	for _, r := range lastMapState {
		if r == '|' || r == '-' || r == '+' || r == '@' {
			result++
		}
	}

	pt1Result = result
	return result, nil
}

func Part2(input string, debug *Debugger) (any, error) {
	slog.Error("this algorithm is invalid currently. It's 82 too high for my input, but the invalids don't make sense so I'm gonna move on...")

	result := 0

	finalState := State{LoopCheck: map[int]int{}}
	for i, r := range input {
		x := i - (finalState.Height * (finalState.Width + 1))

		switch r {
		case '\n':
			if finalState.Width == 0 {
				finalState.Width = i
			}
			finalState.Height++
		case '^', 'V', '<', '>':
			finalState.Guard = Guard{
				Pos: Pos{X: x, Y: finalState.Height},
				Dir: r,
			}
		case '#':
			finalState.Obstacles = append(finalState.Obstacles, Obstacle{Pos: Pos{X: x, Y: finalState.Height}, Symbol: '#'})
		}
	}
	finalState.Height++

	for !finalState.GuardOutOfBounds() {
		finalState = finalState.Step()
	}

	var loopedStates []State
	for i := len(finalState.Paths) - 1; i >= 1; i-- {
		fmt.Printf("\rStep: %d", i)

		// if finalState.Paths[i].Hit.Symbol != 0 {
		// 	continue
		// }

		state := finalState
		state.Obstacles = append(slices.Clone(state.Obstacles), Obstacle{Pos: finalState.Paths[i].Pos, Symbol: 'O'})
		state.Guard = finalState.Paths[i-1]
		state.Paths = nil
		state.LoopCheck = map[int]int{}

		for !state.GuardOutOfBounds() {
			state = state.Step()
			// clearScreen()
			// fmt.Printf("~~%d\n%s", i, state.Debug())

			if state.InLoop {
				// clearScreen()
				// fmt.Printf("~~%d\n%s", i, state.Debug())

				loopedStates = append(loopedStates, state)
				break
			}
		}
	}
	fmt.Print("\r")

	uniqueObstacles := map[string]int{}
	for _, state := range loopedStates {
		o := state.Obstacles[len(state.Obstacles)-1]
		key := fmt.Sprintf("%d,%d", o.Y, o.X)
		if _, ok := uniqueObstacles[key]; !ok {
			result++
			debug.WriteFunc(func() string { return VisualizeStep + key + "\n" + state.Debug() })
		}
		uniqueObstacles[key]++
	}

	return result, nil
}

type Pos struct {
	X, Y int
}

type Obstacle struct {
	Pos
	Symbol rune
}

var zeroObstacle Obstacle

type Guard struct {
	Pos
	Dir rune
	Hit Obstacle
}

func (g Guard) NextPos() Pos {
	p := g.Pos
	switch g.Dir {
	case '^':
		p.Y--
	case 'V':
		p.Y++
	case '<':
		p.X--
	case '>':
		p.X++
	}
	return p
}

func (g Guard) UniqueKey() int {
	// PX  PY  DIR
	// XXX XXX XX
	return g.Y*100_000 + g.X*100 + int(g.Dir)
}

type State struct {
	StepCount int
	Width     int
	Height    int
	Guard     Guard
	Paths     []Guard
	Obstacles []Obstacle

	LoopCheck map[int]int
	InLoop    bool
}

func (s State) GuardOutOfBounds() bool {
	return s.Guard.X < 0 || s.Guard.X >= s.Width || s.Guard.Y < 0 || s.Guard.Y >= s.Height
}

func (s State) Step() State {
	s.StepCount++

	s.Paths = append(s.Paths, s.Guard)
	s.Guard.Hit, s.Guard.Pos = s.CheckHit()
	if s.Guard.Hit.Symbol != 0 {
		s.Guard.Dir = DirectionChange[s.Guard.Dir]

		guardKey := s.Guard.UniqueKey()
		s.LoopCheck[guardKey]++
		if s.LoopCheck[guardKey] > 1 {
			s.InLoop = true
		}
	}

	return s
}

func (s State) CheckHit() (Obstacle, Pos) {
	newPos := s.Guard.NextPos()
	for _, obstacle := range s.Obstacles {
		if newPos.X == obstacle.X && newPos.Y == obstacle.Y {
			return obstacle, s.Guard.Pos
		}
	}
	return zeroObstacle, newPos
}

func (s State) Debug() string {
	var builder strings.Builder
	builder.Grow(s.Height * (s.Width + 1))
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			if s.Guard.X == x && s.Guard.Y == y {
				builder.WriteString(AnsiColorRed)
				builder.WriteRune(s.Guard.Dir)
				builder.WriteString(AnsiColorReset)
				continue
			}

			hasObstacle := false
			for _, obstacle := range s.Obstacles {
				if obstacle.X == x && obstacle.Y == y {
					hasObstacle = true
					if obstacle.Symbol == '#' {
						builder.WriteString(AnsiColorMagenta)
					} else {
						builder.WriteString(AnsiColorBlue)
					}
					builder.WriteRune(obstacle.Symbol)
					builder.WriteString(AnsiColorReset)
					break
				}
			}
			if hasObstacle {
				continue
			}

			pathSymbol := ""
			for i, guard := range s.Paths {
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

func (s State) IsLooped() bool {
	if s.Guard == s.Paths[0] {
		return true
	}

	// hits := map[int]int{}
	// for _, path := range s.Paths {
	// 	key := path.UniqueKey()

	// 	if path.Hit.Symbol != 0 {
	// 		hits[key]++
	// 		if hits[key] > 2 {
	// 			return true
	// 		}
	// 	}

	// 	// hits[key]++
	// 	// if hits[key] > 2 {
	// 	// 	return true
	// 	// }
	// }
	return false
}

func (s State) SpeculateObstacle() (State, bool) {
	if obstacle, _ := s.CheckHit(); obstacle != zeroObstacle {
		return s, false
	}

	s.Paths = slices.Clone(s.Paths)
	s.Obstacles = slices.Clone(s.Obstacles)

	// Put obstacle in front and "hit"
	s.Obstacles = append(s.Obstacles, Obstacle{Pos: s.Guard.NextPos(), Symbol: 'O'})

	s.Paths = append(s.Paths, s.Guard)

	s.Guard.Dir = DirectionChange[s.Guard.Dir]
	s.Guard.Hit = zeroObstacle

	for {
		s = s.Step()

		// fmt.Println(VisualizeStep + s.Debug())

		// Might not work if sim goes farther
		if pt1Result != -1 && s.StepCount > pt1Result*4 {
			return s, false
		}

		if s.GuardOutOfBounds() {
			return s, false
		}

		if s.Guard.Hit.Symbol != 0 {
			hits := map[int]int{}
			for _, path := range s.Paths {
				if path.Hit.Symbol != 0 {
					key := path.UniqueKey()
					hits[key]++
					if hits[key] > 2 {
						return s, true
					}
				}
			}

			uniquePaths := map[int]int{}
			for _, path := range s.Paths {
				// PX  PY  DIR
				// XXX XXX XX
				key := path.Y*1000 + path.X*100 + int(path.Dir)
				uniquePaths[key]++
				if uniquePaths[key] > 1 {
					return s, true
				}
			}
		}
	}
}
