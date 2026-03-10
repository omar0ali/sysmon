package sysmon

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/omar0ali/sysmon/helper"
)

const cpuinfo_path = "/proc/cpuinfo"

type CpuInfo struct {
	LogicalCPUs   int
	PhysicalCores int
	ModelName     string
}

func ReadCpuInfo() CpuInfo {
	var cpuinfo CpuInfo
	helper.OpenScanner(cpuinfo_path, func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue // empty line
			}

			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue // not key and a value
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			switch key {
			case "processor":
				cpuinfo.LogicalCPUs++
			// we only need the following once.
			case "model name":
				if cpuinfo.ModelName == "" {
					cpuinfo.ModelName = value
				}
			case "cpu cores":
				if cpuinfo.PhysicalCores == 0 {
					cpuinfo.PhysicalCores, _ = strconv.Atoi(value)
				}
			}
		}
	})
	return cpuinfo
}
