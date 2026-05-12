package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/anuress/miru/adb"
	"github.com/anuress/miru/picker"
	"github.com/anuress/miru/tui"
)

func main() {
	deviceFlag := flag.String("device", "", "ADB device serial (skip device picker)")
	processFlag := flag.String("process", "", "App package name (skip process picker)")
	portFlag := flag.Int("port", 6360, "OkHttp Profiler port")
	flag.Parse()

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
			m := picker.NewDeviceModel(devices)
			p := tea.NewProgram(m, tea.WithAltScreen())
			result, err := p.Run()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			dm := result.(picker.DeviceModel)
			if dm.Chosen() == nil {
				os.Exit(0)
			}
			serial = dm.Chosen().Serial
		}
	}

	pkg := *processFlag
	if pkg == "" {
		procs, err := adb.Processes(serial)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		m := picker.NewProcessModel(procs)
		p := tea.NewProgram(m, tea.WithAltScreen())
		result, err := p.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		pm := result.(picker.ProcessModel)
		if pm.Chosen() == nil {
			os.Exit(0)
		}
		pkg = pm.Chosen().Package
	}

	if err := adb.Forward(serial, *portFlag); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer adb.RemoveForward(serial, *portFlag)

	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", *portFlag))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not connect to OkHttp Profiler on port %d — is the interceptor running?\n", *portFlag)
		os.Exit(1)
	}
	defer conn.Close()

	app := tui.NewAppModel(conn, serial, pkg, *portFlag)
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
