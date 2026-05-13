package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
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
	cursorLine   int
	lastKey      string
	focused      bool
}

func NewDetailModel() DetailModel {
	return DetailModel{}
}

func (m DetailModel) Init() tea.Cmd { return nil }

func (m *DetailModel) SetRequest(r *model.Request) {
	m.request = r
	m.cursorLine = 0
	m.lastKey = ""
	m.search = NewSearch("")
	m.searching = false
	m.searchInput = ""
	m.currentMatch = 0
}

func (m DetailModel) Update(msg tea.Msg) (DetailModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.searching {
			lines := m.contentLines()
			body := strings.Join(lines, "\n")
			matches := m.search.FindMatches(body)
			switch msg.String() {
			case "esc":
				m.searching = false
				m.searchInput = ""
				m.search = NewSearch("")
				m.currentMatch = 0
			case "up":
				if len(matches) > 0 {
					m.currentMatch = (m.currentMatch - 1 + len(matches)) % len(matches)
					m.cursorLine = matchLine(body, matches[m.currentMatch].Start)
				}
			case "down":
				if len(matches) > 0 {
					m.currentMatch = (m.currentMatch + 1) % len(matches)
					m.cursorLine = matchLine(body, matches[m.currentMatch].Start)
				}
			case "backspace":
				if len(m.searchInput) > 0 {
					m.searchInput = m.searchInput[:len(m.searchInput)-1]
					m.search = NewSearch(m.searchInput)
					m.currentMatch = 0
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
				m.cursorLine = 0
				m.lastKey = ""
			}
		case "right":
			if m.activeTab < tabCount-1 {
				m.activeTab++
				m.cursorLine = 0
				m.lastKey = ""
			}
		case "/":
			m.searching = true
		case "up", "k":
			if m.cursorLine > 0 {
				m.cursorLine--
			}
			m.lastKey = ""
		case "down", "j":
			if m.cursorLine < m.maxCursor() {
				m.cursorLine++
			}
			m.lastKey = ""
		case "ctrl+u":
			half := (m.height - 1) / 2
			m.cursorLine -= half
			if m.cursorLine < 0 {
				m.cursorLine = 0
			}
			m.lastKey = ""
		case "ctrl+d":
			half := (m.height - 1) / 2
			m.cursorLine += half
			if mc := m.maxCursor(); m.cursorLine > mc {
				m.cursorLine = mc
			}
			m.lastKey = ""
		case "G":
			m.cursorLine = m.maxCursor()
			m.lastKey = ""
		case "g":
			if m.lastKey == "g" {
				m.cursorLine = 0
				m.lastKey = ""
			} else {
				m.lastKey = "g"
			}
		default:
			m.lastKey = ""
		}
	}
	return m, nil
}

func (m DetailModel) View() string {
	if m.request == nil {
		return lipgloss.NewStyle().Foreground(ColorGray).Render("\n  select a request to view details")
	}

	// Create tab styles inside View() so lipgloss detects the live terminal context
	activeTab := lipgloss.NewStyle().
		Background(lipgloss.Color("#1f6feb")).
		Foreground(lipgloss.Color("#ffffff")).
		Bold(true).
		PaddingLeft(1).
		PaddingRight(1)
	inactiveTab := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8b949e"))

	var tabs []string
	for i, name := range tabNames {
		if DetailTab(i) == m.activeTab {
			tabs = append(tabs, activeTab.Render(name))
		} else {
			tabs = append(tabs, inactiveTab.Render(name))
		}
	}
	tabBar := strings.Join(tabs, "  ")
	content := m.renderTab()

	// Bound content lines to pane width
	if m.width > 2 {
		content = boundWidth(content, m.width-1)
	}

	// Clip to pane height: tabBar(1) + blank(1) + optional searchBar(1) = 2-3 overhead
	overhead := 2
	if m.searching || m.search.Query != "" {
		overhead = 3
	}
	availH := m.height - overhead
	if availH > 0 {
		content = strings.TrimRight(content, "\n")
		lines := strings.Split(content, "\n")

		// Derive scroll offset to keep cursorLine visible
		scrollOffset := 0
		if m.cursorLine >= availH {
			scrollOffset = m.cursorLine - availH + 1
		}
		start := scrollOffset
		end := start + availH
		if end > len(lines) {
			end = len(lines)
		}

		visible := make([]string, end-start)
		copy(visible, lines[start:end])
		// Only highlight cursor when pane is focused
		if m.focused {
			cursorInView := m.cursorLine - start
			if cursorInView >= 0 && cursorInView < len(visible) {
				plain := stripANSI(visible[cursorInView])
				visible[cursorInView] = "\033[7m" + plain + "\033[0m"
			}
		}
		content = strings.Join(visible, "\n")
	}

	if m.searching || m.search.Query != "" {
		searchBg := lipgloss.NewStyle().
			Background(lipgloss.Color("#1f6feb")).
			Foreground(lipgloss.Color("#ffffff")).
			Width(m.width - 1).
			PaddingLeft(1)
		label := searchBg.Render(fmt.Sprintf("/ %s_   %s   ↑↓ navigate · esc clear",
			m.searchInput, m.matchInfo()))
		return tabBar + "\n\n" + content + "\n" + label
	}
	return tabBar + "\n\n" + content
}

// contentLines returns the rendered, width-bounded, trimmed lines for the current tab.
func (m DetailModel) contentLines() []string {
	if m.request == nil {
		return nil
	}
	content := m.renderTab()
	if m.width > 2 {
		content = boundWidth(content, m.width-1)
	}
	content = strings.TrimRight(content, "\n")
	return strings.Split(content, "\n")
}

// maxCursor returns the last valid cursorLine index.
func (m DetailModel) maxCursor() int {
	lines := m.contentLines()
	if len(lines) == 0 {
		return 0
	}
	return len(lines) - 1
}

func (m DetailModel) respBody() string {
	if m.request == nil {
		return ""
	}
	return m.request.RespBody
}

func (m DetailModel) matchInfo() string {
	matches := m.search.FindMatches(strings.Join(m.contentLines(), "\n"))
	if len(matches) == 0 {
		return "no matches"
	}
	return fmt.Sprintf("%d of %d matches", m.currentMatch+1, len(matches))
}

func (m DetailModel) renderTab() string {
	r := m.request
	var content string
	switch m.activeTab {
	case TabRequestFull:
		content = renderHeaders(r.ReqHeaders) + "\n" + renderBody(r.ReqBody, r.ReqBodyType)
	case TabReqHeaders:
		content = renderHeaders(r.ReqHeaders)
	case TabRespHeaders:
		content = renderHeaders(r.RespHeaders)
	case TabRespBody:
		content = renderBody(r.RespBody, r.RespBodyType)
	default:
		return ""
	}
	if m.search.Query != "" {
		content = highlightMatches(content, m.search, m.currentMatch)
	}
	return content
}

func renderHeaders(headers map[string]string) string {
	if len(headers) == 0 {
		return StyleInactiveTab.Render("  (no headers)")
	}
	keys := make([]string, 0, len(headers))
	for k := range headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(StyleHeaderKey.Render(k) + ": " + headers[k] + "\n")
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
	return label + prettyBody(body)
}

func prettyBody(s string) string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return s
	}
	// Try JSON pretty-print
	if (s[0] == '{' || s[0] == '[') {
		var buf bytes.Buffer
		if err := json.Indent(&buf, []byte(s), "", "  "); err == nil {
			return buf.String()
		}
	}
	return s
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ValueAtCursor returns the smart-copy value for the line under the cursor.
func (m DetailModel) ValueAtCursor() (string, bool) {
	lines := m.contentLines()
	if m.cursorLine >= len(lines) {
		return "", false
	}
	raw := stripANSI(lines[m.cursorLine])
	if isBlockOpener(raw) {
		block := blockCopy(stripANSILines(lines), m.cursorLine)
		if block == "" {
			return "", false
		}
		return block, true
	}
	val := extractValue(raw)
	if val == "" {
		return "", false
	}
	return val, true
}

// LineAtCursor returns the full ANSI-stripped trimmed line under the cursor.
func (m DetailModel) LineAtCursor() (string, bool) {
	lines := m.contentLines()
	if m.cursorLine >= len(lines) {
		return "", false
	}
	raw := strings.TrimSpace(stripANSI(lines[m.cursorLine]))
	if raw == "" {
		return "", false
	}
	return raw, true
}

// matchLine returns the 0-based line index containing byte offset pos in text.
func matchLine(text string, pos int) int {
	if pos > len(text) {
		pos = len(text)
	}
	return strings.Count(text[:pos], "\n")
}

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
