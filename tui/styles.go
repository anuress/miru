package tui

import "github.com/charmbracelet/lipgloss"

// Color vars — set via ApplyTheme() before the TUI starts.
var (
	ColorGreen  lipgloss.Color
	ColorRed    lipgloss.Color
	ColorOrange lipgloss.Color
	ColorBlue   lipgloss.Color
	ColorGray   lipgloss.Color
	ColorWhite  lipgloss.Color
	ColorBg     lipgloss.Color
	ColorBgAlt  lipgloss.Color
	ColorBorder lipgloss.Color
)

// Style functions — created at call time so they always use the current colors.

func StyleInactiveTab() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(ColorGray)
}

func StyleHeaderKey() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(ActiveTheme.HeaderKey)
}

func StyleGrayStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(ColorGray)
}

func StatusStyle(code int, inFlight bool) lipgloss.Style {
	s := lipgloss.NewStyle()
	switch {
	case inFlight:
		return s.Foreground(ColorBlue)
	case code >= 500:
		return s.Foreground(ColorOrange)
	case code >= 400:
		return s.Foreground(ColorRed)
	case code >= 200:
		return s.Foreground(ColorGreen)
	default:
		return s.Foreground(ColorGray)
	}
}

func MethodStyle(method string) lipgloss.Style {
	s := lipgloss.NewStyle().Bold(true)
	switch method {
	case "GET":
		return s.Foreground(ActiveTheme.Blue)
	case "POST":
		return s.Foreground(ActiveTheme.Green)
	case "PUT":
		return s.Foreground(ActiveTheme.Purple)
	case "PATCH":
		return s.Foreground(ActiveTheme.Orange)
	case "DELETE":
		return s.Foreground(ActiveTheme.Red)
	case "HEAD":
		return s.Foreground(ActiveTheme.Gray)
	default:
		return s.Foreground(ActiveTheme.White)
	}
}
