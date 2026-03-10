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
	time.Sleep(1 * time.Second)
	curr := sysmon.ReadCpuStat()

	delta := sysmon.DeltaCPUStats(prev, curr)

	for i, d := range delta {
		total := d.User + d.Nice + d.System + d.Idle + d.Iowait + d.Irq + d.SoftIrq
		idle := d.Idle + d.Iowait
		usage := float64(total-idle) / float64(total) * 100
		fmt.Printf("CPU%d usage: %.1f%%\n", i, usage)
	}
}
