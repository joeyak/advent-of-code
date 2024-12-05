package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joeyak/advent-of-code/common"

	_ "github.com/joeyak/advent-of-code/2024"
)

func main() {
	var fast bool
	var year, day, part int
	flag.BoolVar(&fast, "fast", false, "")
	flag.IntVar(&year, "year", time.Now().Year(), "")
	flag.IntVar(&day, "day", 1, "")
	flag.IntVar(&part, "part", 1, "")
	flag.Parse()

	adventFunc, ok := common.GetAdventFunc(year, day, part)
	if !ok {
		slog.Error("advent func did not exist", "year", year, "day", day, "part", part)
		os.Exit(1)
	}

	input, err := getInput(year, day)
	if err != nil {
		slog.Error("could not get input", "err", err)
		os.Exit(1)
	}

	p := tea.NewProgram(Model{
		f:     adventFunc,
		input: input,
		start: time.Now(),
	})
	if _, err := p.Run(); err != nil {
		slog.Error("program errored", "err", err)
		os.Exit(1)
	}
}

func getInput(year, day int) (string, error) {
	sessionToken, err := os.ReadFile(".session")
	if err != nil {
		return "", fmt.Errorf("could not get session token: %w", err)
	}

	dir := "inputs"
	inputFile := fmt.Sprintf("%s/input-%d-%d.txt", dir, year, day)
	data, err := os.ReadFile(inputFile)
	if err == nil {
		slog.Info("using cached puzzle input")
	} else if errors.Is(err, os.ErrNotExist) {
		slog.Info("downloading puzzle input")

		dirInfo, err := os.Stat(dir)
		if err != nil {
			err := os.Mkdir(dir, 0755)
			if err != nil {
				return "", fmt.Errorf("could not create input dir")
			}
		} else if !dirInfo.IsDir() {
			return "", fmt.Errorf("inputs path is not a directory")
		}

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day), nil)
		if err != nil {
			return "", fmt.Errorf("could not create a new request: %w", err)
		}

		req.AddCookie(&http.Cookie{Name: "session", Value: string(sessionToken)})

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("could not send request: %w", err)
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("could not read body: %w", err)
		}

		err = os.WriteFile(inputFile, data, 0644)
		if err != nil {
			return "", fmt.Errorf("could not write input file: %w", err)
		}
	} else {
		return "", fmt.Errorf("could not read input file: %w", err)
	}

	return string(data), nil
}

type State struct {
	Complete  bool
	Result    any
	Data      any
	Debug     string
	StepCount int
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type Model struct {
	f     common.AdventFunc
	input string
	fast  bool

	start        time.Time
	currentState State
	pastStates   []State
}

func (m Model) step() tea.Msg {
	state, result, debug, err := m.f(m.input, m.currentState)
	complete := errors.Is(err, common.ErrPuzzleComplete)
	if !complete && err != nil {
		return errMsg{err}
	}

	return State{
		Complete:  complete,
		Result:    result,
		StepCount: m.currentState.StepCount + 1,
		Data:      state,
		Debug:     debug,
	}
}

func (m Model) Init() tea.Cmd {
	return m.step
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case " ":
			return m, m.step
		case "left":
			lastIdx := len(m.pastStates) - 1
			m.currentState = m.pastStates[lastIdx]
			m.pastStates = m.pastStates[:lastIdx]
			return m, m.step
		}
	case errMsg:
		slog.Error("an error occured", "err", msg)
		return m, tea.Quit
	case State:
		m.pastStates = append(m.pastStates, m.currentState)
		m.currentState = msg
		if m.fast {
			return m, m.step
		}
	}

	return m, nil
}

func (m Model) View() string {
	statusStyle := lipgloss.NewStyle().Padding(0, 2)
	return lipgloss.JoinVertical(lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Top,
			statusStyle.Render(fmt.Sprintf("Complete: %t", m.currentState.Complete)),
			statusStyle.Render(fmt.Sprintf("Step: %d", m.currentState.StepCount)),
			statusStyle.Render(fmt.Sprintf("Result: %v", m.currentState.Result)),
		),
		m.currentState.Debug,
	)
}
