package common

import (
	"errors"
	"sync"
)

var ErrPuzzleComplete = errors.New("puzzle complete")

type AdventFunc func(input string, state any) (newState, result any, debug string, err error)

var (
	adventFuncsMu sync.Mutex
	adventFuncs   = make(map[int]map[int]map[int]AdventFunc)
)

func RegisterFunc(year, day, part int, part1, part2 AdventFunc) {
	adventFuncsMu.Lock()
	defer adventFuncsMu.Unlock()

	if _, ok := adventFuncs[year]; !ok {
		adventFuncs[year] = map[int]map[int]AdventFunc{}
	}

	if _, ok := adventFuncs[year][day]; !ok {
		adventFuncs[year][day] = map[int]AdventFunc{}
	}

	adventFuncs[year][day][1] = part1
	adventFuncs[year][day][2] = part2
}

func GetAdventFunc(year, day, part int) (AdventFunc, bool) {
	f, ok := adventFuncs[year][day][part]
	return f, ok
}
