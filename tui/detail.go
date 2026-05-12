package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/anuress/miru/model"
)

type DetailTab int

const (
	TabRequestFull DetailTab = iota
	TabReqHeaders
	TabRespHeaders
	TabRespBody
	tabCount
)

var tabNames = []string{"REQUEST (FULL)", "REQ HEADERS", "RESP HEADERS", "RESP BODY"}

type DetailModel struct {
	request      *model.Request
	activeTab    DetailTab
	search       Search
	searching    bool
	searchInput  string
	currentMatch int
	width        int
	height       int
	scrollY      int
}

func NewDetailModel() DetailModel {
	return DetailModel{}
}

func (m DetailModel) Init() tea.Cmd { return nil }

func (m *DetailModel) SetRequest(r *model.Request) {
	m.request = r
	m.scrollY = 0
	m.search = NewSearch("")
	m.searching = false
	m.searchInput = ""
	m.currentMatch = 0
}

func (m DetailModel) Update(msg tea.Msg) (DetailModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.searching {
			switch msg.String() {
			case "esc":
				m.searching = false
				m.searchInput = ""
				m.search = NewSearch("")
				m.currentMatch = 0
			case "backspace":
				if len(m.searchInput) > 0 {
					m.searchInput = m.searchInput[:len(m.searchInput)-1]
					m.search = NewSearch(m.searchInput)
					m.currentMatch = 0
				}
			case "n":
				body := m.respBody()
				matches := m.search.FindMatches(body)
				if len(matches) > 0 {
					m.currentMatch = (m.currentMatch + 1) % len(matches)
				}
			case "N":
				body := m.respBody()
				matches := m.search.FindMatches(body)
				if len(matches) > 0 {
					m.currentMatch = (m.currentMatch - 1 + len(matches)) % len(matches)
				}
			default:
				if len(msg.String()) == 1 {
					m.searchInput += msg.String()
					m.search = NewSearch(m.searchInput)
					m.currentMatch = 0
				}
			}
			return m, nil
		}
		switch msg.String() {
		case "left":
			if m.activeTab > 0 {
				m.activeTab--
				m.scrollY = 0
			}
		case "right":
			if m.activeTab < tabCount-1 {
				m.activeTab++
				m.scrollY = 0
			}
		case "/":
			if m.activeTab == TabRespBody {
				m.searching = true
			}
		case "up", "k":
			if m.scrollY > 0 {
				m.scrollY--
			}
		case "down", "j":
			m.scrollY++
		}
	}
	return m, nil
}

func (m DetailModel) View() string {
	if m.request == nil {
		return lipgloss.NewStyle().Foreground(ColorGray).Render("\n  select a request to view details")
	}

	var tabs []string
	for i, name := range tabNames {
		if DetailTab(i) == m.activeTab {
			tabs = append(tabs, StyleActiveTab.Render(name))
		} else {
			tabs = append(tabs, StyleInactiveTab.Render(name))
		}
	}
	tabBar := strings.Join(tabs, "  ")
	content := m.renderTab()

	// Bound content lines to pane width
	if m.width > 2 {
		content = boundWidth(content, m.width-1)
	}

	if m.searching || m.search.Query != "" {
		searchBar := fmt.Sprintf("/ %s_   %s   n/N: next/prev · esc: dismiss", m.searchInput, m.matchInfo())
		return tabBar + "\n" + content + "\n" + StyleStatusBar.Render(searchBar)
	}
	return tabBar + "\n" + content
}

func (m DetailModel) respBody() string {
	if m.request == nil {
		return ""
	}
	return m.request.RespBody
}

func (m DetailModel) matchInfo() string {
	matches := m.search.FindMatches(m.respBody())
	if len(matches) == 0 {
		return "no matches"
	}
	return fmt.Sprintf("%d of %d matches", m.currentMatch+1, len(matches))
}

func (m DetailModel) renderTab() string {
	r := m.request
	switch m.activeTab {
	case TabRequestFull:
		return renderHeaders(r.ReqHeaders) + "\n" + renderBody(r.ReqBody, r.ReqBodyType)
	case TabReqHeaders:
		return renderHeaders(r.ReqHeaders)
	case TabRespHeaders:
		return renderHeaders(r.RespHeaders)
	case TabRespBody:
		body := renderBody(r.RespBody, r.RespBodyType)
		if m.search.Query != "" {
			body = highlightMatches(body, m.search, m.currentMatch)
		}
		return body
	}
	return ""
}

func renderHeaders(headers map[string]string) string {
	if len(headers) == 0 {
		return StyleInactiveTab.Render("  (no headers)")
	}
	var sb strings.Builder
	for k, v := range headers {
		sb.WriteString(StyleHeaderKey.Render(k) + ": " + v + "\n")
	}
	return sb.String()
}

func renderBody(body, contentType string) string {
	if body == "" {
		return StyleInactiveTab.Render("  (no body)")
	}
	label := ""
	if contentType != "" {
		label = StyleGray.Render(contentType) + "\n"
	}
	return label + body
}

func highlightMatches(text string, s Search, currentMatch int) string {
	matches := s.FindMatches(text)
	if len(matches) == 0 {
		return text
	}
	current := lipgloss.NewStyle().Background(lipgloss.Color("#5a3e1b")).Foreground(ColorOrange)
	other := lipgloss.NewStyle().Background(lipgloss.Color("#2d4a1e")).Foreground(ColorGreen)
	var result strings.Builder
	pos := 0
	for i, m := range matches {
		result.WriteString(text[pos:m.Start])
		style := other
		if i == currentMatch {
			style = current
		}
		result.WriteString(style.Render(text[m.Start:m.End]))
		pos = m.End
	}
	result.WriteString(text[pos:])
	return result.String()
}

var StyleGray = lipgloss.NewStyle().Foreground(ColorGray)

// boundWidth truncates each line in s to at most maxW visible characters,
// stripping ANSI codes when measuring width.
func boundWidth(s string, maxW int) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if visibleLen(line) > maxW {
			lines[i] = truncateLine(line, maxW)
		}
	}
	return strings.Join(lines, "\n")
}

// visibleLen returns the display width of s ignoring ANSI escape sequences.
func visibleLen(s string) int {
	inEsc := false
	n := 0
	for _, r := range s {
		if r == '\033' {
			inEsc = true
			continue
		}
		if inEsc {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEsc = false
			}
			continue
		}
		n++
	}
	return n
}

// truncateLine truncates s to maxW visible chars, appending "…".
func truncateLine(s string, maxW int) string {
	if maxW <= 1 {
		return "…"
	}
	inEsc := false
	n := 0
	for i, r := range s {
		if r == '\033' {
			inEsc = true
			continue
		}
		if inEsc {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEsc = false
			}
			continue
		}
		if n == maxW-1 {
			return s[:i] + "…" + "\033[0m"
		}
		n++
	}
	return s
}
