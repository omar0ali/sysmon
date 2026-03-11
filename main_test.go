package main

import (
	"testing"

	"github.com/omar0ali/sysmon/sysmon"
)

func TestParseStat(t *testing.T) {
	stat := &sysmon.Stat{}
	line := "1105 (dbus-broker-lau) S 1082 1105 1105 0 -1 4194304 789 0 0 0 0 1 0 0 20 0 1 0 1898 8351744 1494"
	sysmon.ParseStatLine(line, stat)
	if stat.PID != 1105 {
		t.Errorf("expected PID 1105, got %d", stat.PID)
	}
	if stat.Comm != "dbus-broker-lau" {
		t.Errorf("expected comm dbus-broker-lau, got %s", stat.Comm)
	}
	if stat.State != "S" {
		t.Errorf("expected state S")
	}
}
