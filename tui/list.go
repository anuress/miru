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
	m.requests = append([]model.Request{r}, m.requests...)
	if m.autoScroll {
		m.cursor = 0
	} else {
		// shift cursor down to keep the same item selected
		m.cursor++
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
			}
			m.autoScroll = (m.cursor == 0)
		case "down", "j":
			visible := m.visible()
			if m.cursor < len(visible)-1 {
				m.cursor++
				m.autoScroll = false
			}
		}
	}
	return m, nil
}

func (m ListModel) View() string {
	// Fixed columns: method(8) + space(1) + status(6) + space(1) + time(7) + space(1) = 24
	// URL column gets the rest, minimum 20
	urlW := 20
	if m.width > 24+urlW {
		urlW = m.width - 24
	}

	gray := lipgloss.NewStyle().Foreground(ColorGray)

	header := gray.Render(fmt.Sprintf("%-8s %-*s %-6s %s", "METHOD", urlW, "URL", "STATUS", "TIME"))

	visible := m.visible()
	var allRows []string
	selStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#1c2d3a")).
		Foreground(lipgloss.Color("#e6edf3"))

	for i, r := range visible {
		url := truncate(r.URL, urlW)
		if i == m.cursor {
			plain := fmt.Sprintf("%-8s %-*s %-6s %s",
				r.Method, urlW, url,
				statusText(r.StatusCode, r.InFlight),
				fmt.Sprintf("%dms", r.Duration.Milliseconds()),
			)
			allRows = append(allRows, selStyle.Width(m.width).Render(plain))
		} else {
			method := MethodStyle(r.Method).Render(fmt.Sprintf("%-8s", r.Method))
			status := StatusStyle(r.StatusCode, r.InFlight).Render(fmt.Sprintf("%-6d", r.StatusCode))
			duration := gray.Render(fmt.Sprintf("%dms", r.Duration.Milliseconds()))
			allRows = append(allRows, fmt.Sprintf("%s %-*s %s %s", method, urlW, url, status, duration))
		}
	}

	// Clip to available height (subtract 1 for header)
	availH := m.height - 1
	if availH <= 0 || len(allRows) <= availH {
		return strings.Join(append([]string{header}, allRows...), "\n")
	}
	offset := m.cursor - availH + 1
	if offset < 0 {
		offset = 0
	}
	end := offset + availH
	if end > len(allRows) {
		end = len(allRows)
	}
	return strings.Join(append([]string{header}, allRows[offset:end]...), "\n")
}

func statusText(code int, inFlight bool) string {
	if inFlight {
		return "..."
	}
	if code == 0 {
		return "???"
	}
	return fmt.Sprintf("%d", code)
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
