package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	"github.com/HuBeZa/automatons/elementary/engine"
)

var (
	gameStyle   = lipgloss.NewStyle().Border(lipgloss.DoubleBorder())
	footerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).MarginLeft(1)
)

type mainModel struct {
	game          engine.Game
	tickWatch     stopwatch.Model
	chanceForLife int
	rule          uint8
}

func newModelCenterBlock(rule uint8, chanceForLife int) tea.Model {
	width, height := getTerminalSize()
	m := mainModel{
		game:          engine.NewCenterBlock(rule, width-2, height-4),
		tickWatch:     stopwatch.NewWithInterval(50 * time.Millisecond),
		chanceForLife: chanceForLife,
		rule:          rule,
	}
	return m
}

func newModelRandom(rule uint8, chanceForLife int) tea.Model {
	width, height := getTerminalSize()
	m := mainModel{
		game:          engine.NewRandom(rule, width-2, height-4, float64(chanceForLife)/100.0),
		tickWatch:     stopwatch.NewWithInterval(50 * time.Millisecond),
		chanceForLife: chanceForLife,
		rule:          rule,
	}
	return m
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			return m, tea.Quit
		case "ctrl+n":
			return newModelRandom(m.rule, m.chanceForLife), nil
		case "alt+n":
			return newModelCenterBlock(m.rule, m.chanceForLife), nil
		case " ":
			return m, m.tickWatch.Toggle()
		case "<":
			m.tickWatch.Interval += 50 * time.Millisecond
		case ">":
			if m.tickWatch.Interval > 0 {
				m.tickWatch.Interval -= 50 * time.Millisecond
			}
		case "[":
			if m.chanceForLife > 0 {
				m.chanceForLife--
			}
		case "]":
			if m.chanceForLife < 100 {
				m.chanceForLife++
			}
		case "{":
			m.rule--
		case "}":
			m.rule++
		}
	case stopwatch.TickMsg:
		if msg.ID == m.tickWatch.ID() {
			var eof bool
			eof, m.game = m.game.Tick()
			if eof {
				return m, m.tickWatch.Stop()
			}
			return m.updateWatch(msg)
		}
	case stopwatch.StartStopMsg, stopwatch.ResetMsg:
		return m.updateWatch(msg)
	}
	return m, nil
}

func (m mainModel) updateWatch(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.tickWatch, cmd = m.tickWatch.Update(msg)
	return m, cmd
}

func (m mainModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		m.renderGame(),
		m.renderFooter())
}

func (m mainModel) renderGame() string {
	var sb strings.Builder
	for row := range m.game.Rows() {
		for col := range m.game.Columns() {
			sb.WriteString(m.renderCell(row, col))
		}
		if row != m.game.Rows()-1 {
			sb.WriteRune('\n')
		}
	}

	return gameStyle.Render(sb.String())
}

func (m mainModel) renderCell(row, col int) string {
	if m.game.Get(row, col) {
		return "█"
	}
	return " "
}

func (m mainModel) renderFooter() string {
	spaceCmd := "play"
	if m.tickWatch.Running() {
		spaceCmd = "pause"
	}

	return footerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			fmt.Sprintf("speed: tick/%v • change: < > • space: %v", m.tickWatch.Interval, spaceCmd),
			fmt.Sprintf("rule: %v • change: { } • chance: %v%% • change: [ ] • ctrl+n: new random • alt+n: new center", m.rule, m.chanceForLife),
		),
	)
}

func main() {
	m := newModelRandom(62, 50)
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func getTerminalSize() (width, height int) {
	width, height, _ = term.GetSize(int(os.Stdout.Fd()))
	return
}
