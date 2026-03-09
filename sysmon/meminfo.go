package sysmon

import (
	"bufio"
	"strings"

	"github.com/omar0ali/activity-monitor/helper"
)

const meminfo_path = "/proc/meminfo"

type MemInfo struct {
	Total     uint64
	Free      uint64
	Available uint64
	Cached    uint64
}

func ReadMemInfo(unit Unit) MemInfo {
	data := map[string]string{}
	helper.OpenScanner(meminfo_path, func(scanner *bufio.Scanner) {
		scanner.Split(bufio.ScanLines) //set lines (default) can be ignored
		for scanner.Scan() {
			parts := strings.SplitN(scanner.Text(), ":", 2)
			if len(parts) != 2 {
				return
			}
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			data[key] = value
		}

	})

	return MemInfo{
		Total:     helper.ParseUint(data["MemTotal"]) / uint64(unit),
		Free:      helper.ParseUint(data["MemFree"]) / uint64(unit),
		Available: helper.ParseUint(data["MemAvailable"]) / uint64(unit),
		Cached:    helper.ParseUint(data["Cached"]) / uint64(unit),
	}
}
