package adb

import (
	"fmt"
	"os/exec"
	"strings"
)

type Device struct {
	Serial string
	Info   string
}

type Process struct {
	PID     string
	Package string
}

func ParseDevices(output string) []Device {
	var devices []Device
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "\tdevice") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				devices = append(devices, Device{Serial: parts[0]})
			}
		}
	}
	return devices
}

func ParseProcesses(output string) []Process {
	var procs []Process
	lines := strings.Split(output, "\n")
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			pkg := fields[len(fields)-1]
			if strings.Contains(pkg, ".") {
				procs = append(procs, Process{PID: fields[1], Package: pkg})
			}
		}
	}
	return procs
}

func Devices() ([]Device, error) {
	out, err := exec.Command("adb", "devices").Output()
	if err != nil {
		return nil, fmt.Errorf("adb not found — install Android SDK platform-tools and ensure adb is on your PATH")
	}
	return ParseDevices(string(out)), nil
}

func Processes(serial string) ([]Process, error) {
	out, err := exec.Command("adb", "-s", serial, "shell", "ps").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list processes on %s: %w", serial, err)
	}
	return ParseProcesses(string(out)), nil
}

func Forward(serial string, port int) error {
	arg := fmt.Sprintf("tcp:%d", port)
	_, err := exec.Command("adb", "-s", serial, "forward", arg, arg).Output()
	if err != nil {
		return fmt.Errorf("adb forward failed: %w", err)
	}
	return nil
}

func RemoveForward(serial string, port int) {
	arg := fmt.Sprintf("tcp:%d", port)
	exec.Command("adb", "-s", serial, "forward", "--remove", arg).Run()
}
