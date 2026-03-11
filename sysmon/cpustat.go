package sysmon

import (
	"bufio"
	"strings"

	"github.com/omar0ali/sysmon/sysmon/helper"
)

const cpustat_path = "/proc/stat"

type CPUStats struct {
	User, Nice, System, Idle, Iowait, Irq, SoftIrq uint64
}

func ReadCpuStat() []CPUStats {
	var cpustat []CPUStats
	helper.OpenScanner(cpustat_path, func(scanner *bufio.Scanner) {
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "cpu") {
				parts := strings.Fields(line)
				if len(parts) < 8 { // void index out of range panic
					continue
				}
				cpustat = append(cpustat, CPUStats{
					User:    helper.ParseUint(parts[1]),
					Nice:    helper.ParseUint(parts[2]),
					System:  helper.ParseUint(parts[3]),
					Idle:    helper.ParseUint(parts[4]),
					Iowait:  helper.ParseUint(parts[5]),
					Irq:     helper.ParseUint(parts[6]),
					SoftIrq: helper.ParseUint(parts[7]),
				})
			}
			if !strings.HasPrefix(line, "cpu") {
				break
			}
		}
	})
	return cpustat
}

func DeltaCPUStats(prev, curr []CPUStats) []CPUStats {
	n := min(len(curr),
		// safety: in case slices are unequal
		len(prev))

	deltas := make([]CPUStats, n)
	for i := 0; i < n; i++ {
		deltas[i] = CPUStats{
			User:    curr[i].User - prev[i].User,
			Nice:    curr[i].Nice - prev[i].Nice,
			System:  curr[i].System - prev[i].System,
			Idle:    curr[i].Idle - prev[i].Idle,
			Iowait:  curr[i].Iowait - prev[i].Iowait,
			Irq:     curr[i].Irq - prev[i].Irq,
			SoftIrq: curr[i].SoftIrq - prev[i].SoftIrq,
		}
	}
	return deltas
}
