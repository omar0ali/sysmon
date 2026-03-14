package main

import (
	"testing"

	"github.com/omar0ali/sysmon/sysmon"
)

func TestParseStatLine(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		wantPID   int
		wantComm  string
		wantState string
	}{
		{
			name:      "normal process",
			line:      "1105 (dbus-broker-lau) S 1082 1105 1105 0 -1 4194304 789 0 0 0 0 1 0 0 20 0 1 0 1898 8351744 1494",
			wantPID:   1105,
			wantComm:  "dbus-broker-lau",
			wantState: "S",
		},
		{
			name:      "zombie process",
			line:      "2345 (my-zombie) Z 1 2345 2345 0 -1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0",
			wantPID:   2345,
			wantComm:  "my-zombie",
			wantState: "Z",
		},
		{
			name:      "process with spaces in name",
			line:      "5678 (my process) R 5677 5678 5678 0 -1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0",
			wantPID:   5678,
			wantComm:  "my process",
			wantState: "R",
		},
		{
			name:      "process with parentheses inside name",
			line:      "91011 (weird(proc)) S 91010 91011 91011 0 -1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0",
			wantPID:   91011,
			wantComm:  "weird(proc)",
			wantState: "S",
		},
		{
			name:      "invalid line missing fields",
			line:      "1234 (short)",
			wantPID:   1234,
			wantComm:  "short",
			wantState: "", // State missing
		},
		{
			name:      "empty line",
			line:      "",
			wantPID:   0,
			wantComm:  "",
			wantState: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stat := &sysmon.Stat{}
			sysmon.ParseStatLine(tt.line, stat)
			if stat.PID != tt.wantPID {
				t.Errorf("expected PID %d, got %d", tt.wantPID, stat.PID)
			}
			if stat.Comm != tt.wantComm {
				t.Errorf("expected comm %q, got %q", tt.wantComm, stat.Comm)
			}
			if stat.State != tt.wantState {
				t.Errorf("expected state %q, got %q", tt.wantState, stat.State)
			}
		})
	}
}
