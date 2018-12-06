package customized

import (
	"reflect"

	"github.com/kricen/pmisc/lib"
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
	Collect() ([]Metric, error, AlarmInfo)
}

type AlarmInfo struct {
	JobName    string `json:"job_name"`
	HostIP     string `json:"host_ip"`
	HostName   string `json:"host_name"`
	MetricName string `json:"metric_name"`
	Reason     string `json:"reason"`
	Timestamp  int64  `json:"timestamp"`
	// has two type system  and global type ,default value is 'global'
	// 'system' : hava maximum attemps,be constrained by alarm gateway
	// 'global' : send alarm every time
	Type string `json:"type"`
}
