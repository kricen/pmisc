package customized

import (
	"reflect"

	"github.com/pmisc/lib"
	"github.com/toolkits/nux"
)

const (
	historyCount int = 2
)

type CpuCollector struct {
	ps *nux.ProcStat
}

func NewCpuCollector() *CpuCollector {

	return &CpuCollector{}
}

func (c *CpuCollector) Collect() (metrics []Metric, err error) {
	c.ps, err = nux.CurrentProcStat()
	if err != nil {
		return
	}
	idle := c.cpuIdle()
	metrics = append(metrics, Metric{Name: "idle", Value: lib.Decimal(idle), NameType: "cpu", ValueType: reflect.Float64, MetricName: "node_cpu"})
	metrics = append(metrics, Metric{Name: "busy", Value: lib.Decimal(1 - idle), NameType: "cpu", ValueType: reflect.Float64, MetricName: "node_cpu"})
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
