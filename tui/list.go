package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/anuress/miru/model"
)

type ListModel struct {
	requests   []model.Request
	cursor     int
	filter     Filter
	autoScroll bool
	width      int
	height     int
}

func NewListModel() ListModel {
	return ListModel{autoScroll: true}
}

func (m *ListModel) AddRequest(r model.Request) {
	m.requests = append(m.requests, r)
	if m.autoScroll {
		m.cursor = len(m.requests) - 1
	}
}

func (m *ListModel) UpdateRequest(r model.Request) {
	for i, req := range m.requests {
		if req.ID == r.ID {
			m.requests[i] = r
			return
		}
	}
}

func (m *ListModel) Clear() {
	m.requests = nil
	m.cursor = 0
}

func (m *ListModel) SetFilter(f Filter) {
	m.filter = f
}

func (m ListModel) Selected() *model.Request {
	visible := m.visible()
	if m.cursor < 0 || m.cursor >= len(visible) {
		return nil
	}
	r := visible[m.cursor]
	return &r
}

func (m ListModel) visible() []model.Request {
	var out []model.Request
	for _, r := range m.requests {
		if m.filter.Match(r) {
			out = append(out, r)
		}
	}
	return out
}

func (m ListModel) Update(msg tea.Msg) (ListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m.autoScroll = false
			}
		case "down", "j":
			visible := m.visible()
			if m.cursor < len(visible)-1 {
				m.cursor++
			}
			if m.cursor == len(visible)-1 {
				m.autoScroll = true
			}
		}
	}
	return m, nil
}

func (m ListModel) View() string {
	header := lipgloss.NewStyle().Foreground(ColorGray).Render(
		fmt.Sprintf("%-8s %-40s %-8s %s", "METHOD", "URL", "STATUS", "TIME"),
	)
	var rows []string
	rows = append(rows, header)
	for i, r := range m.visible() {
		method := StatusStyle(r.StatusCode, r.InFlight).Render(fmt.Sprintf("%-8s", r.Method))
		url := truncate(r.URL, 40)
		status := StatusStyle(r.StatusCode, r.InFlight).Render(fmt.Sprintf("%-8d", r.StatusCode))
		duration := lipgloss.NewStyle().Foreground(ColorGray).Render(
			fmt.Sprintf("%dms", r.Duration.Milliseconds()),
		)
		row := fmt.Sprintf("%s %-40s %s %s", method, url, status, duration)
		if i == m.cursor {
			row = StyleSelected.Render(row)
		}
		rows = append(rows, row)
	}
	return strings.Join(rows, "\n")
}

func truncate(s string, n int) string {
	s = strings.TrimPrefix(s, "https://")
	s = strings.TrimPrefix(s, "http://")
	if idx := strings.Index(s, "/"); idx >= 0 {
		s = s[idx:]
	}
	if len(s) > n {
		return s[:n-1] + "…"
	}
	return fmt.Sprintf("%-*s", n, s)
}
