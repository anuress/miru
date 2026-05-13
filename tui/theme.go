package tui

import "github.com/charmbracelet/lipgloss"

// Theme holds all color tokens for a miru color scheme.
type Theme struct {
	Green     lipgloss.Color
	Red       lipgloss.Color
	Orange    lipgloss.Color
	Blue      lipgloss.Color
	Purple    lipgloss.Color
	Gray      lipgloss.Color
	White     lipgloss.Color
	Bg        lipgloss.Color
	BgAlt     lipgloss.Color
	Border    lipgloss.Color
	Accent    lipgloss.Color // active tabs, active border, search bar bg
	AccentFg  lipgloss.Color // text on accent background
	HeaderKey lipgloss.Color // header key labels in detail pane
}

// ActiveTheme is the currently applied theme. Set via ApplyTheme before the TUI starts.
var ActiveTheme = CatppuccinMocha

// CatppuccinMocha — soothing pastel theme, dark variant.
// Palette: https://github.com/catppuccin/catppuccin
var CatppuccinMocha = Theme{
	Green:     "#a6e3a1", // Green
	Red:       "#f38ba8", // Red
	Orange:    "#fab387", // Peach
	Blue:      "#89b4fa", // Blue
	Purple:    "#cba6f7", // Mauve
	Gray:      "#7f849c", // Overlay1
	White:     "#cdd6f4", // Text
	Bg:        "#1e1e2e", // Base
	BgAlt:     "#181825", // Mantle
	Border:    "#45475a", // Surface1
	Accent:    "#89b4fa", // Blue
	AccentFg:  "#1e1e2e", // Base (dark text on blue)
	HeaderKey: "#74c7ec", // Sapphire
}

// GithubDark — GitHub's dark mode color scheme.
var GithubDark = Theme{
	Green:     "#3fb950",
	Red:       "#f78166",
	Orange:    "#ffa657",
	Blue:      "#58a6ff",
	Purple:    "#d2a8ff",
	Gray:      "#8b949e",
	White:     "#e6edf3",
	Bg:        "#0d1117",
	BgAlt:     "#161b22",
	Border:    "#30363d",
	Accent:    "#1f6feb",
	AccentFg:  "#ffffff",
	HeaderKey: "#79c0ff",
}

// ThemeByName returns the theme for the given name, or false if unknown.
func ThemeByName(name string) (Theme, bool) {
	switch name {
	case "catppuccin", "catppuccin-mocha":
		return CatppuccinMocha, true
	case "github-dark", "github":
		return GithubDark, true
	}
	return Theme{}, false
}

// ApplyTheme updates the active theme and all color vars.
// Must be called before the TUI starts.
func ApplyTheme(t Theme) {
	ActiveTheme = t
	ColorGreen  = t.Green
	ColorRed    = t.Red
	ColorOrange = t.Orange
	ColorBlue   = t.Blue
	ColorGray   = t.Gray
	ColorWhite  = t.White
	ColorBg     = t.Bg
	ColorBgAlt  = t.BgAlt
	ColorBorder = t.Border
}
