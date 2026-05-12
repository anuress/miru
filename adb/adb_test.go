package adb_test

import (
	"testing"

	"github.com/nurliman/miru/adb"
)

func TestParseDevices_MultipleDevices(t *testing.T) {
	input := "List of devices attached\nemulator-5554\tdevice\nR3CT109ABCD\tdevice"
	devices := adb.ParseDevices(input)
	if len(devices) != 2 {
		t.Fatalf("expected 2 devices, got %d", len(devices))
	}
	if devices[0].Serial != "emulator-5554" {
		t.Errorf("expected emulator-5554, got %s", devices[0].Serial)
	}
}

func TestParseDevices_Empty(t *testing.T) {
	input := "List of devices attached\n"
	devices := adb.ParseDevices(input)
	if len(devices) != 0 {
		t.Fatalf("expected 0 devices, got %d", len(devices))
	}
}

func TestParseProcesses_FiltersInput(t *testing.T) {
	input := "USER      PID   PPID  VSZ   RSS   WCHAN  PC        NAME\nu0_a123   12043 1234  1234  1234  0      0         com.myapp.debug\nu0_a124   12044 1234  1234  1234  0      0         com.myapp"
	procs := adb.ParseProcesses(input)
	if len(procs) != 2 {
		t.Fatalf("expected 2 processes, got %d", len(procs))
	}
	if procs[0].Package != "com.myapp.debug" {
		t.Errorf("unexpected package: %s", procs[0].Package)
	}
}
