package customized

import (
	"fmt"
	"reflect"

	"github.com/pmisc/lib"
	"github.com/toolkits/nux"
)

const (
	historyCount int = 2
)

type CpuCollector struct {
	cpuRecord lib.JobQueue
	ps        *nux.ProcStat
}

func NewCpuCollector() *CpuCollector {
	cr := lib.NewJobQueue(150)
	return &CpuCollector{cpuRecord: cr}
}

func (c *CpuCollector) Collect() (metrics []Metric, err error, ai AlarmInfo) {
	c.ps, err = nux.CurrentProcStat()
	if err != nil {
		return
	}
	idle := c.cpuIdle()
	busy := lib.Decimal(1 - idle)
	tmp := c.cpuRecord.Push(busy)
	if tmp != nil {
		var totalUsed float64
		for _, v := range c.cpuRecord.AccessDatas() {
			totalUsed += v.(float64)
		}
		averageUsed := lib.Decimal(totalUsed / float64(c.cpuRecord.AccessLen()))
		if averageUsed >= 0.8 {
			ai.MetricName = "cpu"
			ai.Reason = fmt.Sprintf("cpu使用率超过阈值，近5分钟使用率:%f%s", averageUsed*100, "%")
		}
	}
	metrics = append(metrics, Metric{Name: "idle", Value: lib.Decimal(idle), NameType: "cpu", ValueType: reflect.Float64, MetricName: "node_cpu"})
	metrics = append(metrics, Metric{Name: "busy", Value: busy, NameType: "cpu", ValueType: reflect.Float64, MetricName: "node_cpu"})
	metrics = append(metrics, Metric{Name: "user", Value: lib.Decimal(c.cpuUser()), NameType: "cpu", ValueType: reflect.Float64, MetricName: "node_cpu"})
	metrics = append(metrics, Metric{Name: "nice", Value: lib.Decimal(c.cpuNice()), NameType: "cpu", ValueType: reflect.Float64, MetricName: "node_cpu"})
	metrics = append(metrics, Metric{Name: "system", Value: lib.Decimal(c.cpuSystem()), NameType: "cpu", ValueType: reflect.Float64, MetricName: "node_cpu"})
	metrics = append(metrics, Metric{Name: "iowait", Value: lib.Decimal(c.cpuIowait()), NameType: "cpu", ValueType: reflect.Float64, MetricName: "node_cpu"})
	metrics = append(metrics, Metric{Name: "irq", Value: lib.Decimal(c.cpuIrq()), NameType: "cpu", ValueType: reflect.Float64, MetricName: "node_cpu"})
	metrics = append(metrics, Metric{Name: "softirq", Value: lib.Decimal(c.cpuSoftIrq()), NameType: "cpu", ValueType: reflect.Float64, MetricName: "node_cpu"})
	metrics = append(metrics, Metric{Name: "steal", Value: lib.Decimal(c.cpuSteal()), NameType: "cpu", ValueType: reflect.Float64, MetricName: "node_cpu"})
	metrics = append(metrics, Metric{Name: "guest", Value: lib.Decimal(c.cpuGuest()), NameType: "cpu", ValueType: reflect.Float64, MetricName: "node_cpu"})

	return
}

func (c *CpuCollector) cpuIdle() float64 {
	return float64(c.ps.Cpu.Idle) / float64(c.ps.Cpu.Total)
}

func (c *CpuCollector) cpuUser() float64 {

	return float64(c.ps.Cpu.User) / float64(c.ps.Cpu.Total)
}

func (c *CpuCollector) cpuNice() float64 {

	return float64(c.ps.Cpu.Nice) / float64(c.ps.Cpu.Total)
}

func (c *CpuCollector) cpuSystem() float64 {

	return float64(c.ps.Cpu.System) / float64(c.ps.Cpu.Total)
}

func (c *CpuCollector) cpuIowait() float64 {

	return float64(c.ps.Cpu.Iowait) / float64(c.ps.Cpu.Total)
}

func (c *CpuCollector) cpuIrq() float64 {

	return float64(c.ps.Cpu.Irq) / float64(c.ps.Cpu.Total)
}

func (c *CpuCollector) cpuSoftIrq() float64 {

	return float64(c.ps.Cpu.SoftIrq) / float64(c.ps.Cpu.Total)
}

func (c *CpuCollector) cpuSteal() float64 {

	return float64(c.ps.Cpu.Steal) / float64(c.ps.Cpu.Total)
}

func (c *CpuCollector) cpuGuest() float64 {

	return float64(c.ps.Cpu.Guest) / float64(c.ps.Cpu.Total)
}
