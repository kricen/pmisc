package customized

import (
	"fmt"
	"log"
	"reflect"

	"github.com/kricen/pmisc/lib"
	"github.com/toolkits/nux"
)

type MemoryCollector struct {
	memRecord lib.JobQueue
}

func NewMemoryCollector() *MemoryCollector {
	mr := lib.NewJobQueue(150)
	return &MemoryCollector{memRecord: mr}
}

func (c *MemoryCollector) Collect() (metrics []Metric, err error, ai AlarmInfo) {

	m, err := nux.MemInfo()
	if err != nil {
		log.Println(err)
		return
	}

	memFree := m.MemFree + m.Buffers + m.Cached
	memUsed := m.MemTotal - memFree
	usedPercent := lib.Decimal(float64(memUsed) / float64(m.MemTotal))
	tmp := c.memRecord.Push(usedPercent)
	if tmp != nil {
		var totalUsed float64
		for _, v := range c.memRecord.AccessDatas() {
			totalUsed += v.(float64)
		}
		averageUsed := lib.Decimal(totalUsed / float64(c.memRecord.AccessLen()))
		if averageUsed >= 0.8 {
			ai.MetricName = "memory"
			ai.Reason = fmt.Sprintf("内存使用率超过阈值，近5分钟使用率:%f%s", averageUsed*100, "%")
		}
	}
	metrics = append(metrics, Metric{Name: "used-percent", Value: usedPercent, NameType: "memory", ValueType: reflect.Uint64, MetricName: "node_memory"})
	metrics = append(metrics, Metric{Name: "total", Value: m.MemTotal, NameType: "memory", ValueType: reflect.Uint64, MetricName: "node_memory"})
	metrics = append(metrics, Metric{Name: "used", Value: memUsed, NameType: "memory", ValueType: reflect.Uint64, MetricName: "node_memory"})
	metrics = append(metrics, Metric{Name: "free", Value: memFree, NameType: "memory", ValueType: reflect.Uint64, MetricName: "node_memory"})
	metrics = append(metrics, Metric{Name: "swap_total", Value: m.SwapTotal, NameType: "memory", ValueType: reflect.Uint64, MetricName: "node_memory"})
	metrics = append(metrics, Metric{Name: "swap_used", Value: m.SwapUsed, NameType: "memory", ValueType: reflect.Uint64, MetricName: "node_memory"})
	metrics = append(metrics, Metric{Name: "swap_free", Value: m.SwapFree, NameType: "memory", ValueType: reflect.Uint64, MetricName: "node_memory"})

	return
}
