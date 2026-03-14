package main

import (
	"maps"
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

func TestParseCpuInfoLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		initial  sysmon.CpuInfo
		expected sysmon.CpuInfo
	}{
		{
			name: "empty line",
			line: "",
			initial: sysmon.CpuInfo{
				LogicalCPUs:   1,
				ModelName:     "Intel",
				PhysicalCores: 4,
			},
			expected: sysmon.CpuInfo{
				LogicalCPUs:   1,
				ModelName:     "Intel",
				PhysicalCores: 4,
			},
		},
		{
			name:     "invalid line without colon",
			line:     "processor 0",
			initial:  sysmon.CpuInfo{},
			expected: sysmon.CpuInfo{},
		},
		{
			name:    "processor increments logical cpu",
			line:    "processor : 0",
			initial: sysmon.CpuInfo{},
			expected: sysmon.CpuInfo{
				LogicalCPUs: 1,
			},
		},
		{
			name:    "model name set first time",
			line:    "model name : Intel(R) Core(TM) i7",
			initial: sysmon.CpuInfo{},
			expected: sysmon.CpuInfo{
				ModelName: "Intel(R) Core(TM) i7",
			},
		},
		{
			name: "model name not overwritten",
			line: "model name : AMD Ryzen",
			initial: sysmon.CpuInfo{
				ModelName: "Intel",
			},
			expected: sysmon.CpuInfo{
				ModelName: "Intel",
			},
		},
		{
			name:    "cpu cores parsed",
			line:    "cpu cores : 8",
			initial: sysmon.CpuInfo{},
			expected: sysmon.CpuInfo{
				PhysicalCores: 8,
			},
		},
		{
			name: "cpu cores not overwritten",
			line: "cpu cores : 16",
			initial: sysmon.CpuInfo{
				PhysicalCores: 8,
			},
			expected: sysmon.CpuInfo{
				PhysicalCores: 8,
			},
		},
		{
			name:     "cpu cores invalid number",
			line:     "cpu cores : notanumber",
			initial:  sysmon.CpuInfo{},
			expected: sysmon.CpuInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := tt.initial

			sysmon.ParseCpuInfoLine(tt.line, &cpu)

			if cpu.LogicalCPUs != tt.expected.LogicalCPUs {
				t.Errorf("LogicalCPUs = %d, want %d", cpu.LogicalCPUs, tt.expected.LogicalCPUs)
			}

			if cpu.ModelName != tt.expected.ModelName {
				t.Errorf("ModelName = %q, want %q", cpu.ModelName, tt.expected.ModelName)
			}

			if cpu.PhysicalCores != tt.expected.PhysicalCores {
				t.Errorf("PhysicalCores = %d, want %d", cpu.PhysicalCores, tt.expected.PhysicalCores)
			}
		})
	}
}

func TestParseMemInfoLine(t *testing.T) {
	tests := []struct {
		name  string
		line  string
		start map[string]string
		want  map[string]string
	}{
		{
			name:  "normal line",
			line:  "MemTotal: 16384256 kB",
			start: map[string]string{},
			want: map[string]string{
				"MemTotal": "16384256 kB",
			},
		},
		{
			name:  "trims spaces",
			line:  "MemFree :   12345 kB   ",
			start: map[string]string{},
			want: map[string]string{
				"MemFree": "12345 kB",
			},
		},
		{
			name:  "invalid line without colon",
			line:  "MemTotal 16384256 kB",
			start: map[string]string{},
			want:  map[string]string{},
		},
		{
			name:  "empty value",
			line:  "SwapTotal:",
			start: map[string]string{},
			want: map[string]string{
				"SwapTotal": "",
			},
		},
		{
			name: "overwrites existing key",
			line: "MemTotal: 2000 kB",
			start: map[string]string{
				"MemTotal": "1000 kB",
			},
			want: map[string]string{
				"MemTotal": "2000 kB",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := map[string]string{}

			// copy starting map
			maps.Copy(data, tt.start)

			sysmon.ParseMemInfoLine(tt.line, data)

			if len(data) != len(tt.want) {
				t.Fatalf("map length = %d, want %d", len(data), len(tt.want))
			}

			for k, v := range tt.want {
				if data[k] != v {
					t.Errorf("data[%q] = %q, want %q", k, data[k], v)
				}
			}
		})
	}
}

func TestParseCpuStatLine(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		startLen  int
		wantLen   int
		wantUser  uint64
		wantIdle  uint64
		wantValid bool
	}{
		{
			name:      "normal cpu line",
			line:      "cpu  100 200 300 400 500 600 700",
			startLen:  0,
			wantLen:   1,
			wantUser:  100,
			wantIdle:  400,
			wantValid: true,
		},
		{
			name:      "cpu0 line",
			line:      "cpu0 10 20 30 40 50 60 70",
			startLen:  0,
			wantLen:   1,
			wantUser:  10,
			wantIdle:  40,
			wantValid: true,
		},
		{
			name:      "not cpu line",
			line:      "intr 123 456",
			startLen:  0,
			wantLen:   0,
			wantValid: false,
		},
		{
			name:      "too few fields",
			line:      "cpu  1 2 3",
			startLen:  0,
			wantLen:   0,
			wantValid: true,
		},
		{
			name:      "empty line",
			line:      "",
			startLen:  0,
			wantLen:   0,
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := make([]*sysmon.CPUStats, tt.startLen)

			ok := sysmon.ParseCpuStatLine(tt.line, &stats)

			if ok != tt.wantValid {
				t.Errorf("return = %v, want %v", ok, tt.wantValid)
			}

			if len(stats) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(stats), tt.wantLen)
			}

			if tt.wantLen > 0 {
				if stats[0].User != tt.wantUser {
					t.Errorf("User = %d, want %d", stats[0].User, tt.wantUser)
				}
				if stats[0].Idle != tt.wantIdle {
					t.Errorf("Idle = %d, want %d", stats[0].Idle, tt.wantIdle)
				}
			}
		})
	}
}
