package main

import (
	"fmt"

	"github.com/pmisc/prometheus/customized"
	"github.com/pmisc/registry"
)

func main() {

	// // dc := customized.DiskCollector{}
	// cc := customized.CpuCollector{}
	// fmt.Println(reflect.TypeOf(cc) == reflect.TypeOf(cc))
	// mc := customized.MemoryCollector{}
	// for i := 0; i < 10; i++ {
	// 	fmt.Println("--------------", i, "--------------")
	// 	metrics, err := mc.Collect()
	// 	if err != nil {
	// 		fmt.Printf("error:%s", err.Error())
	// 	}
	// 	fmt.Println(len(metrics))
	// 	for _, metric := range metrics {
	// 		fmt.Println(metric.Name, metric.Value, metric.NameType, metric.ValueType, metric.MetricName)
	// 	}
	// }
	// var wg sync.WaitGroup
	// rc := &customized.RequestCollector{}
	// wg.Add(2)
	// go func() {
	// 	defer wg.Done()
	// 	for i := 0; i < 10; i++ {
	// 		time.Sleep(1 * time.Second)
	// 		rc.AddRecord("hello", int64(i))
	// 	}
	// }()
	// go func() {
	// 	defer wg.Done()
	// 	for i := 0; i < 2; i++ {
	// 		time.Sleep(2 * time.Second)
	// 		metrics, err := rc.Collect()
	// 		if err != nil {
	// 			fmt.Printf("error:%s", err.Error())
	// 		}
	// 		for _, metric := range metrics {
	// 			fmt.Println(metric.Name, metric.Value, metric.NameType, metric.ValueType, metric.MetricName)
	// 		}
	// 	}
	// }()
	// wg.Wait()
	// fmt.Println("--------------", "--------------")
	// metrics, err := rc.Collect()
	// if err != nil {
	// 	fmt.Printf("error:%s", err.Error())
	// }
	// for _, metric := range metrics {
	// 	fmt.Println(metric.Name, metric.Value, metric.NameType, metric.ValueType, metric.MetricName)
	// }

	cc := customized.NewCpuCollector()
	dc := customized.NewDiskCollector()
	mc := customized.NewMemoryCollector()
	rc := customized.NewRequestCollector()

	cr := registry.NewCollectorRegister("monitor-helper-test")
	cr.Register(cc)
	cr.Register(dc)
	cr.Register(mc)
	cr.Register(rc)
	fmt.Println(cr.ToString())
}
