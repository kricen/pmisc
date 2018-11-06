package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/pmisc/prometheus/customized"
)

func main() {

	// dc := customized.DiskCollector{}
	// cc := customized.CpuCollector{}
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
	var wg sync.WaitGroup
	rc := &customized.RequestCollector{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			rc.AddRecord("hello", int64(i))
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 2; i++ {
			time.Sleep(2 * time.Second)
			metrics, err := rc.Collect()
			if err != nil {
				fmt.Printf("error:%s", err.Error())
			}
			for _, metric := range metrics {
				fmt.Println(metric.Name, metric.Value, metric.NameType, metric.ValueType, metric.MetricName)
			}
		}
	}()
	wg.Wait()
	fmt.Println("--------------", "--------------")
	metrics, err := rc.Collect()
	if err != nil {
		fmt.Printf("error:%s", err.Error())
	}
	for _, metric := range metrics {
		fmt.Println(metric.Name, metric.Value, metric.NameType, metric.ValueType, metric.MetricName)
	}
}
