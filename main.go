package main

import (
	"fmt"
	"time"

	"github.com/omar0ali/sysmon/sysmon"
)

func main() {

	// MemInfo And CpuInfo

	meminfo := sysmon.ReadMemInfo(sysmon.MB)
	fmt.Printf("MEMINFO: %+v\n", meminfo)
	cpuinfo := sysmon.ReadCpuInfo()
	fmt.Printf("CPUINFO: %+v\n", cpuinfo)

	// CpuStat

	prev := sysmon.ReadCpuStat()
	for i, stat := range prev {
		fmt.Printf("CPUSTATS_PREV %d: %+v\n", i, *stat)
	}

	time.Sleep(time.Second)
	curr := sysmon.ReadCpuStat()
	for i, stat := range curr {
		fmt.Printf("CPUSTATS_CURR %d: %+v\n", i, *stat)
	}

	delta := sysmon.DeltaCPUStats(prev, curr)

	for i, d := range delta {
		total := d.User + d.Nice + d.System + d.Idle + d.Iowait + d.Irq + d.SoftIrq
		idle := d.Idle + d.Iowait
		usage := float64(total-idle) / float64(total) * 100
		fmt.Printf("CPU%d usage: %.1f%%\n", i, usage)
	}

	// display processes

	pids := sysmon.GetPids()
	fmt.Printf("Pids: %+v\n", pids)

	procs := map[int]*sysmon.Process{}

	for i := range pids {
		procs[pids[i]] = sysmon.NewProcess(pids[i])
	}

	for i := range procs {
		fmt.Printf("PROC: %+v\n", procs[i])
	}

	// refresh processes to get processes usage

	prevCPU := sysmon.ReadCpuStat()
	prevProcCPU := map[int]uint64{}

	for pid, p := range procs {
		prevProcCPU[pid] = p.Stat.UTime + p.Stat.STime
	}

	time.Sleep(time.Second)
	sysmon.RefreshProcesses(procs)

	currCPU := sysmon.ReadCpuStat()
	delta = sysmon.DeltaCPUStats(prevCPU, currCPU)

	totalDelta := uint64(0)

	for _, d := range delta {
		totalDelta += d.User + d.Nice + d.System + d.Idle +
			d.Iowait + d.Irq + d.SoftIrq
	}

	for pid, p := range procs {

		prev, ok := prevProcCPU[pid]
		if !ok {
			continue
		}

		curr := p.Stat.UTime + p.Stat.STime
		procDelta := curr - prev

		cpuPercent := float64(procDelta) / float64(totalDelta) * 100

		name := p.Stat.Comm

		fmt.Printf("Name: %s PID: %d CPU: %.2f%%\n", name, pid, cpuPercent)
	}
}
