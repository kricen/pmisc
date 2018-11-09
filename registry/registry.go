package registry

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/pmisc/lib"
	"github.com/pmisc/prometheus/customized"
)

type CollectorRegister struct {
	Endpoint   string
	JobName    string
	HostName   string
	URL        string
	collectors []customized.ICollector
}

var (
	DefaultHostIP string
	errorLogger   *log.Logger
)

func init() {
	DefaultHostIP = lib.ResolveHostIP()
	errorLogger = log.New(os.Stdout,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func NewCollectorRegister(jobname string, URL string) *CollectorRegister {
	cr := &CollectorRegister{}
	cr.JobName = jobname
	cr.Endpoint = DefaultHostIP
	cr.URL = URL
	hostname, err := os.Hostname()
	if err != nil {
		hostname = DefaultHostIP
	}
	cr.HostName = hostname
	return cr
}

func (cr *CollectorRegister) Registe(c customized.ICollector) {
	// check whether has same collector
	for _, tc := range cr.collectors {
		if reflect.TypeOf(tc) == reflect.TypeOf(c) {
			return
		}
	}
	cr.collectors = append(cr.collectors, c)
}

func (cr *CollectorRegister) Collect() (metrics []customized.Metric) {

	for _, collector := range cr.collectors {
		ms, err := collector.Collect()
		if err != nil {
			continue
		}
		metrics = append(metrics, ms...)
	}

	return
}

func (cr *CollectorRegister) Push() error {
	// package request url http://localhost:9091/metrics/job/%s/instance/%s
	reqURL := fmt.Sprintf("%s/metrics/job/%s/instance/%s/hostname/%s", cr.URL, cr.JobName, cr.Endpoint, cr.HostName)
	metrics := cr.Collect()
	var ms string
	distinctMetrics := make(map[string]int, 0)
	for _, metric := range metrics {
		gatherName := fmt.Sprintf("%s%s", metric.MetricName, metric.Name)

		count := distinctMetrics[gatherName]
		if count > 0 {
			continue
		}
		distinctMetrics[gatherName]++
		ms += fmt.Sprintf("%s{label=\"%s\"} %+v\n", metric.MetricName, metric.Name, metric.Value)
	}
	req, err := http.NewRequest("POST", reqURL, bytes.NewBufferString(ms))
	if err != nil {
		return err
	}
	resp, err := lib.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if string(bs) != "" {
		return errors.New(string(bs))
	}
	return nil
}

func (cr *CollectorRegister) cornTask() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				err := cr.Push()
				if err != nil {
					errorLogger.Println("something is wrrong,err reason:", err.Error())
				}
			}
		}
	}()
}

func (cr *CollectorRegister) ToString() string {
	return fmt.Sprintf("Endpoint:%s,JobName:%s,collectors size:%d", cr.Endpoint, cr.JobName, len(cr.collectors))
}
