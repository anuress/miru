package picker

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/anuress/miru/adb"
)

func pickerMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// PickDevice shows a Bubble Tea device picker and returns the chosen serial.
// Returns "" if the user quits.
func PickDevice(devices []adb.Device) string {
	m := NewDevicePickerModel(devices)
	p := tea.NewProgram(m, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		return ""
	}
	dm := result.(DevicePickerModel)
	if dm.Chosen() == nil {
		return ""
	}
	return dm.Chosen().Serial
}

// PickProcess shows a Bubble Tea process picker and returns the chosen package.
// Returns "" if the user quits.
func PickProcess(procs []adb.Process, device string) string {
	m := NewProcessPickerModel(procs, device)
	p := tea.NewProgram(m, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		return ""
	}
	pm := result.(ProcessPickerModel)
	if pm.Chosen() == nil {
		return ""
	}
	return pm.Chosen().Package
}

// numbered fallbacks for when Bubble Tea is unavailable
func pickDeviceNumbered(devices []adb.Device) string {
	fmt.Println("Select device:")
	for i, d := range devices {
		fmt.Printf("  %d. %s\n", i+1, d.Serial)
	}
	fmt.Print("Enter number: ")
	var n int
	fmt.Scan(&n)
	if n < 1 || n > len(devices) {
		return ""
	}
	return devices[n-1].Serial
}

func pickProcessNumbered(procs []adb.Process) string {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Filter process (or number to select): ")
		if !scanner.Scan() {
			return ""
		}
		input := strings.TrimSpace(scanner.Text())
		var matched []adb.Process
		for _, p := range procs {
			if strings.Contains(strings.ToLower(p.Package), strings.ToLower(input)) {
				matched = append(matched, p)
			}
		}
		if len(matched) == 1 {
			return matched[0].Package
		}
		for i, p := range matched {
			fmt.Printf("  %d. %s\n", i+1, p.Package)
		}
	}
}
