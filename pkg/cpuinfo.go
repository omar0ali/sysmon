package pkg

import (
	"bufio"
	"github.com/omar0ali/sysmon/internal"
	"strconv"
	"strings"
)

const cpuinfoPath = internal.PROC_DIR + "/cpuinfo"

type CpuInfo struct {
	LogicalCPUs   int
	PhysicalCores int
	ModelName     string
}

func ParseCpuInfoLine(line string, cpuinfo *CpuInfo) {
	if line == "" {
		return
	}

	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	switch key {

	case "processor":
		cpuinfo.LogicalCPUs++

	case "model name":
		if cpuinfo.ModelName == "" {
			cpuinfo.ModelName = value
		}

	case "cpu cores":
		// fallback if physical/core ids don't exist
		if cpuinfo.PhysicalCores == 0 {
			if cores, err := strconv.Atoi(value); err == nil {
				cpuinfo.PhysicalCores = cores
			}
		}
	}
}

func ReadCpuInfo() (*CpuInfo, error) {
	cpuinfo := &CpuInfo{}
	err := internal.OpenWithScanner(cpuinfoPath, func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			ParseCpuInfoLine(line, cpuinfo)
		}
	})

	return cpuinfo, err
}
