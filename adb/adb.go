package adb

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Device struct {
	Serial string
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

// LogcatSession wraps an adb logcat subprocess.
type LogcatSession struct {
	cmd *exec.Cmd
	rc  io.ReadCloser
}

func (s *LogcatSession) Read(p []byte) (int, error) {
	return s.rc.Read(p)
}

// Stop kills the adb logcat subprocess.
func (s *LogcatSession) Stop() {
	if s.cmd.Process != nil {
		s.cmd.Process.Kill()
	}
	s.rc.Close()
	s.cmd.Wait()
}

// StartLogcat starts "adb -s <serial> logcat" filtered by pid if non-empty.
// Returns a LogcatSession whose Read streams raw logcat output.
func StartLogcat(serial, pid string) (*LogcatSession, error) {
	args := []string{"-s", serial, "logcat"}
	if pid != "" {
		args = append(args, "--pid="+pid)
	}
	cmd := exec.Command("adb", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("logcat pipe failed: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start adb logcat: %w", err)
	}
	return &LogcatSession{cmd: cmd, rc: stdout}, nil
}
