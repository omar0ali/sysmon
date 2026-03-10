package sysmon

import (
	"bufio"
	"strings"

	"github.com/omar0ali/sysmon/helper"
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
