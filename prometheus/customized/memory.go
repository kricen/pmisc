package customized

import (
	"log"
	"reflect"

	"github.com/toolkits/nux"
)

type MemoryCollector struct {
}

func NewMemoryCollector() *MemoryCollector {

	return &MemoryCollector{}
}

func (c *MemoryCollector) Collect() (metrics []Metric, err error) {

	m, err := nux.MemInfo()
	if err != nil {
		log.Println(err)
		return
	}

	memFree := m.MemFree + m.Buffers + m.Cached
	memUsed := m.MemTotal - memFree

	metrics = append(metrics, Metric{Name: "total", Value: m.MemTotal, NameType: "memory", ValueType: reflect.Uint64, MetricName: "node-memory"})
	metrics = append(metrics, Metric{Name: "used", Value: memUsed, NameType: "memory", ValueType: reflect.Uint64, MetricName: "node-memory"})
	metrics = append(metrics, Metric{Name: "free", Value: memFree, NameType: "memory", ValueType: reflect.Uint64, MetricName: "node-memory"})
	metrics = append(metrics, Metric{Name: "swap-total", Value: m.SwapTotal, NameType: "memory", ValueType: reflect.Uint64, MetricName: "node-memory"})
	metrics = append(metrics, Metric{Name: "swap-used", Value: m.SwapUsed, NameType: "memory", ValueType: reflect.Uint64, MetricName: "node-memory"})
	metrics = append(metrics, Metric{Name: "swap-free", Value: m.SwapFree, NameType: "memory", ValueType: reflect.Uint64, MetricName: "node-memory"})

	return
}
