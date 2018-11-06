package customized

import (
	"reflect"

	"github.com/pmisc/lib"
)

var (
	DefaultHostIP = "127.0.0.1"
)

type Metric struct {
	Name       string       `json:"name"`
	Value      interface{}  `json:"value"`
	NameType   string       `json:"name_type"`
	ValueType  reflect.Kind `json:"value_type"`
	MetricName string       `json:"metric_name"`
	Alert      bool         `json:"alert"`
}

type MetricPack struct {
	Endpoint string   `json:"instance"`
	JobName  string   `json:"job"`
	metrics  []Metric `json:"metrics"`
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
