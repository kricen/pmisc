package customized

import (
	"math"
	"reflect"
	"sync"
)

type RequestCollector struct {
	globalMutex sync.RWMutex
	rgMap       map[string]int64 //request guage map
	rdMap       map[string][]int64
}

func NewRequestCollector() *RequestCollector {

	return &RequestCollector{}
}

func (c *RequestCollector) Collect() (metrics []Metric, err error) {
	c.globalMutex.Lock()
	defer c.globalMutex.Unlock()
	if c.rgMap == nil || c.rdMap == nil {
		c.initMap()
	}
	for url, total := range c.rgMap {
		metrics = append(metrics, Metric{Name: url, Value: total, NameType: "request_total", ValueType: reflect.Int64, MetricName: "job_request_total"})
	}

	for url, durations := range c.rdMap {
		duration := averageDuration(durations)
		metrics = append(metrics, Metric{Name: url, Value: duration, NameType: "request_duration", ValueType: reflect.Int64, MetricName: "job_request_duration"})
		c.rdMap[url] = make([]int64, 0)
	}

	return
}

func (c *RequestCollector) AddRecord(url string, msDuration int64) {
	c.globalMutex.Lock()
	defer c.globalMutex.Unlock()
	if c.rgMap == nil || c.rdMap == nil {
		c.initMap()
	}
	if c.rgMap[url] > math.MaxInt64-1 {
		c.rgMap[url] = 0
	}
	c.rgMap[url]++

	c.rdMap[url] = append(c.rdMap[url], msDuration)

}
func averageDuration(durations []int64) (duration int64) {
	if len(durations) == 0 {
		return
	}
	for _, value := range durations {
		duration += value
	}
	duration = duration / int64(len(durations))
	return
}

func (c *RequestCollector) initMap() {
	if c.rdMap == nil {
		c.rdMap = make(map[string][]int64, 0)
	}
	if c.rgMap == nil {
		c.rgMap = make(map[string]int64, 0)
	}
}
