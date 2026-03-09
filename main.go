package main

import (
	"fmt"

	"github.com/omar0ali/activity-monitor/sysmon"
)

func main() {
	meminfo := sysmon.ReadMemInfo(sysmon.MB)
	fmt.Printf("%+v", meminfo)
}
