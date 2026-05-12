package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.design/x/clipboard"

	"github.com/anuress/miru/model"
)

type CurlOverlay struct {
	visible bool
	curl    string
}

func (o *CurlOverlay) Show(r model.Request) {
	o.curl = GenerateCurl(r)
	o.visible = true
}

func (o *CurlOverlay) Hide() {
	o.visible = false
}

func (o CurlOverlay) Update(msg tea.Msg) (CurlOverlay, tea.Cmd) {
	if !o.visible {
		return o, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			clipboard.Write(clipboard.FmtText, []byte(o.curl))
			o.visible = false
		case "esc":
			o.visible = false
		}
	}
	return o, nil
}

func (o CurlOverlay) View() string {
	if !o.visible {
		return ""
	}
	content := o.curl + "\n\n" +
		lipgloss.NewStyle().Foreground(ColorGray).Render("↵ copy to clipboard · esc dismiss")
	return StyleOverlay.Render(content)
}
