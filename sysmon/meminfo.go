package sysmon

import (
	"bufio"
	"strings"

	"github.com/omar0ali/sysmon/sysmon/helper"
)

const meminfo_path = helper.PROC_DIR + "/meminfo"

type MemInfo struct {
	Total     uint64
	Free      uint64
	Available uint64
	Cached    uint64
	Buffers   uint64
	SwapTotal uint64
	SwapFree  uint64
}

func ParseMemInfoLine(line string, data map[string]string) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	data[key] = value
}

func ReadMemInfo(unit Unit) (*MemInfo, error) {
	data := map[string]string{}
	err := helper.OpenWithScanner(meminfo_path, func(scanner *bufio.Scanner) {
		scanner.Split(bufio.ScanLines) //set lines (default) can be ignored
		for scanner.Scan() {
			ParseMemInfoLine(scanner.Text(), data)
		}
	})

	if err != nil {
		return nil, err
	}

	return &MemInfo{
		Total:     helper.ParseUint(data["MemTotal"]) / uint64(unit),
		Free:      helper.ParseUint(data["MemFree"]) / uint64(unit),
		Available: helper.ParseUint(data["MemAvailable"]) / uint64(unit),
		Cached:    helper.ParseUint(data["Cached"]) / uint64(unit),
		Buffers:   helper.ParseUint(data["Buffers"]) / uint64(unit),
		SwapTotal: helper.ParseUint(data["SwapTotal"]) / uint64(unit),
		SwapFree:  helper.ParseUint(data["SwapFree"]) / uint64(unit),
	}, nil
}
