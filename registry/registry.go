package registry

import (
	"fmt"
	"reflect"

	"github.com/pmisc/lib"
	"github.com/pmisc/prometheus/customized"
)

type CollectorRegister struct {
	mp         customized.MetricPack
	collectors []customized.ICollector
}

var DefaultHostIP string

func init() {
	DefaultHostIP = lib.ResolveHostIP()
}

func NewCollectorRegister(jobname string) *CollectorRegister {
	cr := &CollectorRegister{}
	cr.mp.JobName = jobname
	cr.mp.Endpoint = DefaultHostIP
	return cr
}

func (cr *CollectorRegister) Register(c customized.ICollector) {
	// check whether has same collector
	for _, tc := range cr.collectors {
		if reflect.TypeOf(tc) == reflect.TypeOf(c) {
			return
		}
	}
	cr.collectors = append(cr.collectors, c)
}

func (cr *CollectorRegister) ToString() string {
	return fmt.Sprintf("Endpoint:%s,JobName:%s,collectors size:%d", cr.mp.Endpoint, cr.mp.JobName, len(cr.collectors))
}
