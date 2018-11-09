package customized

import (
	"reflect"

	"github.com/pmisc/lib"
)

var (
	DefaultHostIP = "127.0.0.1"
)

// MetricName{Name:,Endpoint:,JobName:,HostName} Value
type Metric struct {
	Name       string       `json:"name"`  // monitor metric
	Value      interface{}  `json:"value"` //metric decimal value
	NameType   string       `json:"name_type"`
	ValueType  reflect.Kind `json:"value_type"`  // metric value's type
	MetricName string       `json:"metric_name"` // metric name
}

func init() {
	// access localhost ip dynimicly
	localhostIP := lib.ResolveHostIP()
	if localhostIP != "" {
		DefaultHostIP = localhostIP
	}
}

type ICollector interface {
	// collect metrics
	Collect() ([]Metric, error)
}
