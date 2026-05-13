package picker

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/anuress/miru/adb"
)

type DevicePickerModel struct {
	devices []adb.Device
	cursor  int
	chosen  *adb.Device
	width   int
	height  int
}

func NewDevicePickerModel(devices []adb.Device) DevicePickerModel {
	return DevicePickerModel{devices: devices}
}

func (m DevicePickerModel) Chosen() *adb.Device { return m.chosen }

func (m DevicePickerModel) Init() tea.Cmd { return nil }

func (m DevicePickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.cursor < len(m.devices)-1 {
				m.cursor++
			}
		case "enter":
			d := m.devices[m.cursor]
			m.chosen = &d
			return m, tea.Quit
		case "esc", "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m DevicePickerModel) View() string {
	if m.width == 0 {
		return ""
	}

	blue := lipgloss.Color("#58a6ff")
	gray := lipgloss.Color("#8b949e")
	white := lipgloss.Color("#e6edf3")
	selectedBg := lipgloss.Color("#1c2d3a")
	bgAlt := lipgloss.Color("#161b22")

	titleStyle := lipgloss.NewStyle().Foreground(blue).Bold(true)
	grayStyle := lipgloss.NewStyle().Foreground(gray)
	whiteStyle := lipgloss.NewStyle().Foreground(white)
	selectedStyle := lipgloss.NewStyle().Background(selectedBg).Foreground(blue).Bold(true)

	boxW := m.width - 4
	if boxW > 70 {
		boxW = 70
	}
	innerW := boxW - 2

	title := titleStyle.Render("◆ miru") + grayStyle.Render("  select device")
	sep := grayStyle.Render(strings.Repeat("─", innerW-2))

	var rows []string
	for i, d := range m.devices {
		label := fmt.Sprintf("%-*s", innerW-4, d.Serial)
		if i == m.cursor {
			rows = append(rows, selectedStyle.Render("  ▶ "+label))
		} else {
			rows = append(rows, whiteStyle.Render("    "+label))
		}
	}

	inner := strings.Join([]string{title, "", sep}, "\n") + "\n" + strings.Join(rows, "\n")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(blue).
		Background(bgAlt).
		Width(innerW).
		Render(inner)

	footer := grayStyle.Render("  ↑↓  navigate  ·  ↵  select  ·  q  quit")

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
