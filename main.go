package main

// this main package aim at to help you how to use it

import (
	"fmt"
	"time"

	"github.com/pmisc/prometheus/customized"
	"github.com/pmisc/registry"
)

func main() {

	cc := customized.NewCpuCollector()
	dc := customized.NewDiskCollector()
	mc := customized.NewMemoryCollector()
	rc := customized.NewRequestCollector()

	cr := registry.NewCollectorRegister("monitor-helper-test", "http://localhost:9091")
	cr.Registe(cc)
	cr.Registe(dc)
	cr.Registe(mc)
	cr.Registe(rc)
	fmt.Println(cr.ToString())
	go func() {
		for i := 0; i < 100; i++ {
			time.Sleep(100 * time.Millisecond)
			rc.AddRecord("getUserName", int64(i))
		}
	}()

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		cr.Push()
	}

}
