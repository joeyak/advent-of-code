package main

import (
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/crazy3lf/colorconv"
)

const (
	VisualizeStep = "==========STEP==========\n"
	VisualizeData = "==========DATA==========\n"
	VisualizeEnd  = "==========END==========\n"
)

func main() {
	var args struct {
		Year             int           `arg:"positional"`
		Day              int           `arg:"positional"`
		Part             int           `arg:"positional" default:"1"`
		AutoPlay         bool          `arg:"-a"`
		DebugHeat        bool          `arg:"-h"`
		AutoPlayDuration time.Duration `arg:"-d" default:"500ms"`
	}
	arg.MustParse(&args)
	// args.Part = 1

	model := Model{
		CurrentStep: 1,
		StepMod:     1,
		Paused:      !args.AutoPlay,
		DebugHeat:   args.DebugHeat,
	}

	debugData, err := os.ReadFile(fmt.Sprintf("../%d/day%d/debug-Part%d.txt", args.Year, args.Day, args.Part))
	if err != nil {
		slog.Error("could not read debug file", "err", err)
		os.Exit(1)
	}

	debug := strings.ReplaceAll(string(debugData), "\r\n", "\n")
	if !strings.Contains(debug, VisualizeStep) {
		slog.Error("debug file has no steps")
		os.Exit(1)
	}

	for _, s := range strings.Split(strings.Split(debug, VisualizeEnd)[0], VisualizeStep)[1:] {
		s := strings.TrimSuffix(s, "\n")
		data := StepData{Data: s}
		if strings.Contains(s, VisualizeData) {
			parts := strings.Split(s, VisualizeData)
			data.Meta = strings.TrimSuffix(parts[0], "\n")
			data.Data = strings.TrimSuffix(parts[1], "\n")
		}
		model.Steps = append(model.Steps, data)
	}

	p := tea.NewProgram(model)

	go func() {
		for {
			p.Send(new(Tick))
			time.Sleep(args.AutoPlayDuration)
		}
	}()

	if _, err := p.Run(); err != nil {
		slog.Error("there has been an error", "err", err)
		os.Exit(1)
	}
}

type Tick tea.Msg

type StepData struct {
	Data string
	Meta string
}

type Model struct {
	Paused      bool
	DebugHeat   bool
	StepMod     int
	CurrentStep int

	Steps []StepData
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	foreward := func() {
		maxStep := len(m.Steps)
		m.CurrentStep = min(m.CurrentStep+m.StepMod, maxStep)
		if m.CurrentStep == maxStep {
			m.Paused = true
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "left":
			m.CurrentStep = max(m.CurrentStep-m.StepMod, 1)
		case "right":
			foreward()
		case "ctrl+left":
			m.CurrentStep = 1
		case "ctrl+right":
			m.CurrentStep = len(m.Steps)
		case "alt+1":
			m.StepMod = 1
		case "alt+2":
			m.StepMod = 2
		case "alt+3":
			m.StepMod = 3
		case "alt+4":
			m.StepMod = 4
		case "alt+5":
			m.StepMod = 5
		case "alt+6":
			m.StepMod = 6
		case "alt+7":
			m.StepMod = 7
		case "alt+8":
			m.StepMod = 8
		case "alt+9":
			m.StepMod = 9
		case " ":
			m.Paused = !m.Paused
		case "h":
			m.DebugHeat = !m.DebugHeat
		case "up":
			if m.StepMod == 1 {
				m.StepMod = 5
			} else if m.StepMod == 5 {
				m.StepMod = 10
			} else {
				m.StepMod += 10
			}
		case "down":
			if m.StepMod == 10 {
				m.StepMod = 5
			} else if m.StepMod == 5 {
				m.StepMod = 1
			} else {
				m.StepMod -= 10
			}
			if m.StepMod < 1 {
				m.StepMod = 1
			}
		}
	case Tick:
		if !m.Paused {
			foreward()
		}
	}

	return m, nil
}

func (m Model) View() string {
	headerStyle := lipgloss.NewStyle().Margin(0, 2)
	header := lipgloss.JoinHorizontal(lipgloss.Center,
		headerStyle.Render(fmt.Sprintf("Step(1-%d): %d", len(m.Steps), m.CurrentStep)),
		headerStyle.Render(fmt.Sprintf("Paused: %t", m.Paused)),
		headerStyle.Render(fmt.Sprintf("Step Mod: %d", m.StepMod)),
	)

	currentStep := m.Steps[m.CurrentStep-1]
	current := []byte(currentStep.Data)
	heatmap := slices.Repeat([]int{0}, len(current))

	maxHeat := 0
	steps := m.Steps[max(m.CurrentStep-max(20, m.StepMod), 0):min(m.CurrentStep-1, len(m.Steps))]
	for _, step := range steps {
		for heatIdx, r := range current {
			if r != step.Data[heatIdx] {
				heatmap[heatIdx]++
				maxHeat = max(heatmap[heatIdx], maxHeat)
			}
		}
	}

	for i, heat := range heatmap {
		if heat > 0 {
			heatmap[i] = maxHeat - heat + 1
		}
	}

	result := ""
	heatDebug := ""
	for i, b := range current {
		if b == '\n' {
			heatDebug += "\n"
			result += "\n"
			continue
		}

		hex := ""
		heat := heatmap[i]
		if heat > 0 {
			color, err := colorconv.HSVToColor(float64(min(280, (heat-1)*20)), 1, 1)
			if err != nil {
				slog.Error("could not convert heatmap to color", "heat", heat, "err", err)
			}
			hex = strings.Replace(colorconv.ColorToHex(color), "0x", "#", 1)
		}

		result += lipgloss.NewStyle().
			Foreground(lipgloss.Color(hex)).
			Render(string(b))
		heatDebug += fmt.Sprintf(" %02d", heatmap[i])
	}

	if m.DebugHeat {
		result = lipgloss.JoinHorizontal(lipgloss.Center, result, heatDebug)
	} else if currentStep.Meta != "" {
		result = lipgloss.JoinHorizontal(lipgloss.Top, result, lipgloss.NewStyle().Padding(0, 1).Render(currentStep.Meta))
	}

	return lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		Padding(0, 1).
		MarginBottom(1).
		Render(lipgloss.JoinVertical(lipgloss.Left, header, result))
}
