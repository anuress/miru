package picker

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/anuress/miru/adb"
)

type DeviceModel struct {
	devices []adb.Device
	cursor  int
	chosen  *adb.Device
}

func NewDeviceModel(devices []adb.Device) DeviceModel {
	return DeviceModel{devices: devices}
}

func (m DeviceModel) Chosen() *adb.Device {
	return m.chosen
}

func (m DeviceModel) Init() tea.Cmd { return nil }

func (m DeviceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m DeviceModel) View() string {
	s := "◆ miru — select device\n\n"
	for i, d := range m.devices {
		cursor := "  "
		if i == m.cursor {
			cursor = "▶ "
		}
		s += fmt.Sprintf("%s%s\n", cursor, d.Serial)
	}
	s += "\n↑↓ navigate · ↵ select · q quit"
	return s
}
