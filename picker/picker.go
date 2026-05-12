package picker

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/anuress/miru/adb"
)

const (
	clrScreen  = "\033[H\033[2J"
	clrLine    = "\r\033[K"
	cursorUp   = "\033[A"
	styleReset = "\033[0m"
	styleSel   = "\033[7m\033[1m" // reverse + bold
)

// PickDevice shows an interactive device list and returns the chosen serial.
// Returns "" if the user quits.
func PickDevice(devices []adb.Device) string {
	cursor := 0
	redraw := func() {
		h := termHeight()
		availH := h - 5 // header=3 lines, footer=2 lines
		if availH < 1 {
			availH = 1
		}
		offset := 0
		if cursor >= availH {
			offset = cursor - availH + 1
		}
		end := offset + availH
		if end > len(devices) {
			end = len(devices)
		}

		fmt.Print(clrScreen)
		fmt.Print("◆ miru — select device\r\n\r\n")
		for i := offset; i < end; i++ {
			d := devices[i]
			if i == cursor {
				fmt.Printf("  %s▶ %-40s%s\r\n", styleSel, d.Serial, styleReset)
			} else {
				fmt.Printf("    %-40s\r\n", d.Serial)
			}
		}
		fmt.Print("\r\n↑↓ / j k  navigate · enter  select · q  quit\r\n")
	}

	if err := setRaw(true); err != nil {
		// raw mode unavailable — fall back to numbered list
		return pickDeviceNumbered(devices)
	}
	defer setRaw(false)

	redraw()
	buf := make([]byte, 4)
	for {
		n, _ := os.Stdin.Read(buf)
		key := buf[:n]
		switch {
		case isUp(key):
			if cursor > 0 {
				cursor--
			}
		case isDown(key):
			if cursor < len(devices)-1 {
				cursor++
			}
		case isEnter(key):
			setRaw(false)
			fmt.Print(clrScreen)
			return devices[cursor].Serial
		case isQuit(key):
			setRaw(false)
			fmt.Print(clrScreen)
			return ""
		}
		redraw()
	}
}

// PickProcess shows an interactive filterable process list and returns the chosen package.
// Returns "" if the user quits.
func PickProcess(procs []adb.Process, device string) string {
	cursor := 0
	filter := ""

	visible := func() []adb.Process {
		var out []adb.Process
		for _, p := range procs {
			if strings.Contains(strings.ToLower(p.Package), strings.ToLower(filter)) {
				out = append(out, p)
			}
		}
		return out
	}

	redraw := func(vis []adb.Process) {
		h := termHeight()
		// header=4 lines, footer=2 lines → available for list rows
		availH := h - 6
		if availH < 1 {
			availH = 1
		}

		// scroll offset: keep cursor within the visible window
		offset := 0
		if cursor >= availH {
			offset = cursor - availH + 1
		}
		end := offset + availH
		if end > len(vis) {
			end = len(vis)
		}

		fmt.Print(clrScreen)
		fmt.Printf("◆ miru — select process  [%s]\r\n\r\n", device)
		fmt.Printf("  Filter: %s_\r\n  %d processes\r\n\r\n", filter, len(vis))
		for i := offset; i < end; i++ {
			p := vis[i]
			if i == cursor {
				fmt.Printf("  %s▶ %-50s (PID %-6s)%s\r\n", styleSel, p.Package, p.PID, styleReset)
			} else {
				fmt.Printf("    %-50s (PID %-6s)\r\n", p.Package, p.PID)
			}
		}
		fmt.Print("\r\n↑↓ / j k  navigate · type  filter · backspace  erase · enter  select · q  quit\r\n")
	}

	if err := setRaw(true); err != nil {
		return pickProcessNumbered(procs)
	}
	defer setRaw(false)

	vis := visible()
	redraw(vis)

	buf := make([]byte, 4)
	for {
		n, _ := os.Stdin.Read(buf)
		key := buf[:n]
		vis = visible()
		if cursor >= len(vis) {
			cursor = max(0, len(vis)-1)
		}
		switch {
		case isUp(key):
			if cursor > 0 {
				cursor--
			}
		case isDown(key):
			if cursor < len(vis)-1 {
				cursor++
			}
		case isEnter(key):
			if len(vis) > 0 {
				chosen := vis[cursor].Package
				setRaw(false)
				fmt.Print(clrScreen)
				return chosen
			}
		case isBackspace(key):
			if len(filter) > 0 {
				filter = filter[:len(filter)-1]
				cursor = 0
			}
		case isQuit(key):
			setRaw(false)
			fmt.Print(clrScreen)
			return ""
		case isPrintable(key):
			filter += string(key[:1])
			cursor = 0
		}
		vis = visible()
		redraw(vis)
	}
}

func isUp(b []byte) bool {
	return string(b) == "k" || string(b) == "\033[A"
}
func isDown(b []byte) bool {
	return string(b) == "j" || string(b) == "\033[B"
}
func isEnter(b []byte) bool {
	return len(b) > 0 && (b[0] == '\r' || b[0] == '\n')
}
func isQuit(b []byte) bool {
	return string(b) == "q" || (len(b) > 0 && b[0] == 3) // ctrl+c
}
func isBackspace(b []byte) bool {
	return len(b) > 0 && (b[0] == 127 || b[0] == 8)
}
func isPrintable(b []byte) bool {
	return len(b) == 1 && b[0] >= 32 && b[0] < 127
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// numbered fallbacks for when raw mode is unavailable
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
