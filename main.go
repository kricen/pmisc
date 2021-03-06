package main

import (
	"fmt"
	"time"

	"github.com/kricen/pmisc/prometheus/customized"
	"github.com/kricen/pmisc/registry"
)

// this main package aim at to help you how to use it

func main() {

	cc := customized.NewCpuCollector()
	dc := customized.NewDiskCollector()
	mc := customized.NewMemoryCollector()
	rc := customized.NewRequestCollector()

	cr := registry.NewCollectorRegister("monitor-helper-test", []string{"http://47.104.159.222:2379", "http://47.104.3.204:2379"})
	cr.Registe(cc)
	cr.Registe(dc)
	cr.Registe(mc)
	cr.Registe(rc)
	fmt.Println(cr.ToString())
	// 发送报警信息，
	cr.SendAlarm("customized", "something is error", "global")
	go func() {
		for i := 0; i < 1000000; i++ {
			time.Sleep(100 * time.Millisecond)
			// add a request record with a specific
			rc.AddRecord("getUserName", int64(i))
		}
	}()

	for {
		time.Sleep(100 * time.Second)
		// cr.Push()
	}

}
