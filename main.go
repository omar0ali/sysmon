package main

import (
	"fmt"

	"github.com/omar0ali/sysmon/sysmon"
)

func main() {
	meminfo := sysmon.ReadMemInfo(sysmon.MB)
	fmt.Printf("%+v\n", meminfo)
	cpuinfo := sysmon.ReadCpuInfo()
	fmt.Printf("%+v\n", cpuinfo)
	cpustat := sysmon.ReadCpuStat()
	fmt.Printf("%+v\n", cpustat)
}
