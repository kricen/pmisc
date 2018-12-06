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

	cr := registry.NewCollectorRegister("monitor-helper-test", []string{"http://localhost:2379", "http://localhost1:2379"})
	cr.Registe(cc)
	cr.Registe(dc)
	cr.Registe(mc)
	cr.Registe(rc)
	fmt.Println(cr.ToString())
	// send alarm message
	cr.SendAlarm("customized", "something is error", "global")
	go func() {
		for i := 0; i < 100; i++ {
			time.Sleep(100 * time.Millisecond)
			rc.AddRecord("getUserName", int64(i))
		}
	}()

	for {
		time.Sleep(10 * time.Second)
		// cr.Push()
	}

}
