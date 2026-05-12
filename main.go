package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/anuress/miru/adb"
	"github.com/anuress/miru/picker"
	"github.com/anuress/miru/tui"
)

func main() {
	deviceFlag := flag.String("device", "", "ADB device serial (skip device picker)")
	processFlag := flag.String("process", "", "App package name (skip process picker)")
	flag.Parse()

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
