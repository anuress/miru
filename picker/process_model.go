package picker

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/anuress/miru/adb"
)

type ProcessPickerModel struct {
	all     []adb.Process
	visible []adb.Process
	cursor  int
	filter  string
	chosen  *adb.Process
	device  string
	width   int
	height  int
}

func NewProcessPickerModel(procs []adb.Process, device string) ProcessPickerModel {
	m := ProcessPickerModel{all: procs, device: device}
	m.applyFilter()
	return m
}

func (m ProcessPickerModel) Chosen() *adb.Process { return m.chosen }

func (m *ProcessPickerModel) applyFilter() {
	m.visible = nil
	for _, p := range m.all {
		if strings.Contains(strings.ToLower(p.Package), strings.ToLower(m.filter)) {
			m.visible = append(m.visible, p)
		}
	}
	if m.cursor >= len(m.visible) {
		m.cursor = pickerMax(0, len(m.visible)-1)
	}
}

func (m ProcessPickerModel) Init() tea.Cmd { return nil }

func (m ProcessPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
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
				m.cursor = 0
				m.applyFilter()
			}
		case "esc", "ctrl+c", "q":
			return m, tea.Quit
		default:
			if len(msg.String()) == 1 {
				m.filter += msg.String()
				m.cursor = 0
				m.applyFilter()
			}
		}
	}
	return m, nil
}

func (m ProcessPickerModel) View() string {
	if m.width == 0 {
		return ""
	}

	// Styles — created inside View() for proper terminal color detection
	blue := lipgloss.Color("#58a6ff")
	gray := lipgloss.Color("#8b949e")
	white := lipgloss.Color("#e6edf3")
	selectedBg := lipgloss.Color("#1c2d3a")
	bgAlt := lipgloss.Color("#161b22")

	titleStyle := lipgloss.NewStyle().Foreground(blue).Bold(true)
	grayStyle := lipgloss.NewStyle().Foreground(gray)
	whiteStyle := lipgloss.NewStyle().Foreground(white)
	selectedStyle := lipgloss.NewStyle().Background(selectedBg).Foreground(blue).Bold(true)

	// Box width — capped at 90 or terminal width
	boxW := m.width - 4
	if boxW > 90 {
		boxW = 90
	}
	innerW := boxW - 2 // border takes 2

	// Title bar
	deviceLabel := grayStyle.Render("  " + m.device)
	title := "  " + titleStyle.Render("◆ miru") + deviceLabel

	// Filter row
	count := grayStyle.Render(fmt.Sprintf("%d / %d", len(m.visible), len(m.all)))
	filterText := whiteStyle.Render(m.filter) + grayStyle.Render("_")
	filterLabel := grayStyle.Render("  Filter: ")
	padding := strings.Repeat(" ", pickerMax(0, innerW-lipgloss.Width(filterLabel)-lipgloss.Width(filterText)-lipgloss.Width(count)-2))
	filterRow := filterLabel + filterText + padding + count

	// Separator
	sep := grayStyle.Render(strings.Repeat("─", innerW-2))

	// Available rows for the list
	// box overhead: border(2) + title(1) + blank(1) + filter(1) + sep(1) + blank(1) + footer(1) = 8
	listH := m.height - 12
	if listH < 3 {
		listH = 3
	}

	offset := 0
	if m.cursor >= listH {
		offset = m.cursor - listH + 1
	}
	end := offset + listH
	if end > len(m.visible) {
		end = len(m.visible)
	}

	var rows []string
	pkgW := innerW - 16 // leave room for PID column
	for i := offset; i < end; i++ {
		p := m.visible[i]
		pidStr := grayStyle.Render(fmt.Sprintf("PID %s", p.PID))
		pkg := p.Package
		if len(pkg) > pkgW {
			pkg = pkg[:pkgW-1] + "…"
		}
		if i == m.cursor {
			sel := selectedStyle.Render(fmt.Sprintf("  ▶ %-*s", pkgW, pkg))
			rows = append(rows, sel+"  "+pidStr)
		} else {
			rows = append(rows, whiteStyle.Render(fmt.Sprintf("    %-*s", pkgW, pkg))+"  "+pidStr)
		}
	}

	// Assemble inner content
	inner := strings.Join([]string{
		title,
		"",
		filterRow,
		sep,
	}, "\n")
	if len(rows) > 0 {
		inner += "\n" + strings.Join(rows, "\n")
	} else {
		inner += "\n" + grayStyle.Render("  no matches")
	}

	// Wrap in rounded box
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(blue).
		Background(bgAlt).
		Width(innerW).
		Render(inner)

	// Footer outside the box
	footer := grayStyle.Render("  ↑↓  navigate  ·  type to filter  ·  ↵  select  ·  q  quit")

	// Center horizontally and vertically
	boxLines := strings.Split(box, "\n")
	leftPad := pickerMax(0, (m.width-lipgloss.Width(box))/2)
	topPad := pickerMax(0, (m.height-len(boxLines)-2)/2)

	var out strings.Builder
	out.WriteString(strings.Repeat("\n", topPad))
	pad := strings.Repeat(" ", leftPad)
	for _, line := range boxLines {
		out.WriteString(pad + line + "\n")
	}
	out.WriteString(pad + footer)
	return out.String()
}
