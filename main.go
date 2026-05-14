package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/anuress/miru/adb"
	"github.com/anuress/miru/picker"
	"github.com/anuress/miru/tui"
)

type miruConfig struct {
	Theme string `json:"theme"`
}

func loadConfig() miruConfig {
	home, err := os.UserHomeDir()
	if err != nil {
		return miruConfig{Theme: "catppuccin-mocha"}
	}
	path := filepath.Join(home, ".config", "miru", "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return miruConfig{Theme: "catppuccin-mocha"}
	}
	var c miruConfig
	if err := json.Unmarshal(data, &c); err != nil {
		return miruConfig{Theme: "catppuccin-mocha"}
	}
	return c
}

func main() {
	deviceFlag := flag.String("device", "", "ADB device serial (skip device picker)")
	processFlag := flag.String("process", "", "App package name (skip process picker)")
	themeFlag := flag.String("theme", "", "Color theme: catppuccin-mocha, github-dark (overrides config file)")
	flag.Parse()

	// Apply theme: flag > config file > default (catppuccin-mocha)
	cfg := loadConfig()
	themeName := cfg.Theme
	if *themeFlag != "" {
		themeName = *themeFlag
	}
	if t, ok := tui.ThemeByName(themeName); ok {
		tui.ApplyTheme(t)
	} else {
		tui.ApplyTheme(tui.CatppuccinMocha)
	}

	// Resolve device
	serial := *deviceFlag
	if serial == "" {
		devices, err := adb.Devices()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if len(devices) == 0 {
			fmt.Fprintln(os.Stderr, "no devices found — connect a device or start an emulator")
			os.Exit(1)
		}
		if len(devices) == 1 {
			serial = devices[0].Serial
		} else {
			serial = picker.PickDevice(devices)
			if serial == "" {
				os.Exit(0)
			}
		}
	}

	// Resolve process + PID
	pkg := *processFlag
	pid := ""
	if pkg == "" {
		procs, err := adb.Processes(serial)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		pkg = picker.PickProcess(procs, serial)
		if pkg == "" {
			os.Exit(0)
		}
		// find the PID for the chosen package
		for _, p := range procs {
			if p.Package == pkg {
				pid = p.PID
				break
			}
		}
	} else {
		// --process flag given, look up PID
		procs, err := adb.Processes(serial)
		if err == nil {
			for _, p := range procs {
				if p.Package == pkg {
					pid = p.PID
					break
				}
			}
		}
	}

	// Start logcat session
	session, err := adb.StartLogcat(serial, pid)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer session.Stop()

	// Launch TUI
	app := tui.NewAppModel(session, serial, pkg, pid)
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
