package tui

import "github.com/charmbracelet/lipgloss"

var (
	ColorGreen  = lipgloss.Color("#3fb950")
	ColorRed    = lipgloss.Color("#f78166")
	ColorOrange = lipgloss.Color("#ffa657")
	ColorBlue   = lipgloss.Color("#58a6ff")
	ColorGray   = lipgloss.Color("#8b949e")
	ColorWhite  = lipgloss.Color("#e6edf3")
	ColorBg     = lipgloss.Color("#0d1117")
	ColorBgAlt  = lipgloss.Color("#161b22")
	ColorBorder = lipgloss.Color("#30363d")

	StyleActiveTab = lipgloss.NewStyle().
			Foreground(ColorBlue).
			BorderBottom(true).
			BorderForeground(ColorBlue)

	StyleInactiveTab = lipgloss.NewStyle().
				Foreground(ColorGray)

	StyleSelected = lipgloss.NewStyle().
			Background(lipgloss.Color("#1c2d3a")).
			BorderLeft(true).
			BorderForeground(ColorBlue)

	StyleHeaderKey = lipgloss.NewStyle().Foreground(lipgloss.Color("#79c0ff"))
	StyleString    = lipgloss.NewStyle().Foreground(lipgloss.Color("#a5d6ff"))
	StyleNumber    = lipgloss.NewStyle().Foreground(ColorOrange)

	StyleStatusBar = lipgloss.NewStyle().
			Background(ColorBgAlt).
			Foreground(ColorGray)

	StyleOverlay = lipgloss.NewStyle().
			Background(ColorBgAlt).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBlue).
			Padding(1, 2)
)

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
		return s.Foreground(lipgloss.Color("#58a6ff")) // blue
	case "POST":
		return s.Foreground(lipgloss.Color("#3fb950")) // green
	case "PUT":
		return s.Foreground(lipgloss.Color("#d2a8ff")) // purple
	case "PATCH":
		return s.Foreground(lipgloss.Color("#ffa657")) // orange
	case "DELETE":
		return s.Foreground(lipgloss.Color("#f78166")) // red
	case "HEAD":
		return s.Foreground(lipgloss.Color("#8b949e")) // gray
	default:
		return s.Foreground(lipgloss.Color("#e6edf3")) // white
	}
}
