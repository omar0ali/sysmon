package sysmon

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/omar0ali/sysmon/sysmon/helper"
)

type Stat struct {
	PID        int    // field 1
	Comm       string // field 2
	State      string // field 3
	PPID       int    // field 4
	UTime      uint64 // field 14
	STime      uint64 // field 15
	Priority   int    // field 18
	Nice       int    // field 19
	NumThreads int    // field 20
	StartTime  uint64 // field 22
	VSize      uint64 // field 23
	RSS        int64  // field 24
}

type Status struct {
	Name    string
	UID     int
	VmSize  uint64
	VmRSS   uint64
	Threads int
}

type Cmdline struct {
	Args []string
}

type Process struct {
	PID int

	Stat    Stat
	Status  Status
	Cmdline Cmdline
}

func ParseStatLine(line string, stat *Stat) {
	start := strings.Index(line, "(")
	end := strings.LastIndex(line, ")")
	stat.PID, _ = strconv.Atoi(strings.TrimSpace(line[:start]))
	stat.Comm = line[start+1 : end]
	fields := strings.Fields(line[end+2:])
	stat.State = fields[0]
	stat.PPID, _ = strconv.Atoi(fields[1])
	stat.UTime, _ = strconv.ParseUint(fields[11], 10, 64)
	stat.STime, _ = strconv.ParseUint(fields[12], 10, 64)
	stat.Priority, _ = strconv.Atoi(fields[15])
	stat.Nice, _ = strconv.Atoi(fields[16])
	stat.NumThreads, _ = strconv.Atoi(fields[17])
	stat.StartTime, _ = strconv.ParseUint(fields[19], 10, 64)
	stat.VSize, _ = strconv.ParseUint(fields[20], 10, 64)
	stat.RSS, _ = strconv.ParseInt(fields[21], 10, 64)
}

func ParseStatusLine(line string, status *Status) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return
	}

	key := parts[0]
	value := strings.TrimSpace(parts[1])

	switch key {

	case "Name":
		status.Name = value

	case "Uid":
		fields := strings.Fields(value)
		uid, _ := strconv.Atoi(fields[0])
		status.UID = uid

	case "VmSize":
		fields := strings.Fields(value)
		v, _ := strconv.ParseUint(fields[0], 10, 64)
		status.VmSize = v

	case "VmRSS":
		fields := strings.Fields(value)
		v, _ := strconv.ParseUint(fields[0], 10, 64)
		status.VmRSS = v

	case "Threads":
		v, _ := strconv.Atoi(value)
		status.Threads = v
	}
}

func GetPids() ([]int, error) {
	var pids []int
	dir, err := os.ReadDir(helper.PROC_DIR)
	if err != nil {
		return nil, err
	}

	for i := range dir {
		if dir[i].IsDir() {
			v, err := strconv.Atoi(dir[i].Name())
			if err != nil {
				log.Printf("Ignore file, is not pid: %s\n", err)
				continue
			}
			pids = append(pids, v)
		}
	}
	return pids, nil
}

func ParseStat(pid int) (*Stat, error) {
	stat := &Stat{}
	const file = "stat"
	path := fmt.Sprintf("%s/%d/%s", helper.PROC_DIR, pid, file)
	err := helper.OpenWithScanner(path, func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			ParseStatLine(line, stat)
		}

	})

	if err != nil {
		return nil, err
	}
	return stat, nil
}

func ParseCmdline(pid int) *Cmdline {
	const file = "cmdline"
	path := fmt.Sprintf("%s/%d/%s", helper.PROC_DIR, pid, file)

	data, err := os.ReadFile(path)
	if err != nil {
		return &Cmdline{Args: []string{}}
	}
	parts := strings.Split(string(data), "\x00")
	// remove trailing empty entry
	if len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}

	return &Cmdline{
		Args: parts,
	}
}

func ParseStatus(pid int) *Status {
	const file = "status"
	path := fmt.Sprintf("%s/%d/%s", helper.PROC_DIR, pid, file)
	status := &Status{}
	helper.OpenWithScanner(path, func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			line := scanner.Text()
			ParseStatusLine(line, status)
		}
	})

	return status
}

func NewProcess(pid int) (*Process, error) {
	stat, err := ParseStat(pid)
	if err != nil {
		return nil, err
	}
	status := ParseStatus(pid)
	cmdline := ParseCmdline(pid)

	return &Process{
		PID:     pid,
		Stat:    *stat,
		Status:  *status,
		Cmdline: *cmdline,
	}, nil
}

func (p *Process) updateStat() error {
	stat, err := ParseStat(p.PID) // reads /proc/[pid]/stat
	if err != nil {
		return err
	}
	p.Stat.UTime = stat.UTime
	p.Stat.STime = stat.STime
	p.Stat.StartTime = stat.StartTime
	return nil
}

// used to get an update of the UTime, STime, StartTime after sleep by i.e one second
// must read all processes first

func RefreshProcesses(procs map[int]*Process) error {
	pids, err := GetPids()
	if err != nil {
		return err
	}
	seen := map[int]bool{}
	for _, pid := range pids {
		seen[pid] = true
		if p, ok := procs[pid]; ok {
			p.updateStat()
		} else {
			procs[pid], err = NewProcess(pid)
			if err != nil {
				return err
			}
		}
	}

	for pid := range procs {
		if !seen[pid] {
			delete(procs, pid)
		}
	}
	return nil
}
