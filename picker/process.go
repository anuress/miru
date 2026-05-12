package picker

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/anuress/miru/adb"
)

const (
	ansiReset   = "\033[0m"
	ansiReverse = "\033[7m"
	ansiBold    = "\033[1m"
)

type ProcessModel struct {
	all     []adb.Process
	visible []adb.Process
	cursor  int
	filter  string
	chosen  *adb.Process
}

func NewProcessModel(procs []adb.Process) ProcessModel {
	m := ProcessModel{all: procs}
	m.applyFilter()
	return m
}

func (m ProcessModel) Chosen() *adb.Process {
	return m.chosen
}

func (m *ProcessModel) applyFilter() {
	m.visible = nil
	for _, p := range m.all {
		if strings.Contains(strings.ToLower(p.Package), strings.ToLower(m.filter)) {
			m.visible = append(m.visible, p)
		}
	}
	m.cursor = 0
}

func (m ProcessModel) Init() tea.Cmd { return nil }

func (m ProcessModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.visible)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.visible) > 0 {
				p := m.visible[m.cursor]
				m.chosen = &p
				return m, tea.Quit
			}
		case "backspace":
			if len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.applyFilter()
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			if len(msg.String()) == 1 {
				m.filter += msg.String()
				m.applyFilter()
			}
		}
	}
	return m, nil
}

func (m ProcessModel) View() string {
	s := fmt.Sprintf("◆ miru — select process\n\nFilter: %s_\n%d processes\n\n", m.filter, len(m.visible))
	for i, p := range m.visible {
		line := fmt.Sprintf("  %s  (PID %s)", p.Package, p.PID)
		if i == m.cursor {
			line = ansiReverse + ansiBold + fmt.Sprintf("▶ %s  (PID %s)", p.Package, p.PID) + ansiReset
		}
		s += line + "\n"
	}
	s += "\n↑↓ navigate · type to filter · ↵ select · q quit"
	return s
}
