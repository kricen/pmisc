package registry

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/kricen/pmisc/lib"
	"github.com/kricen/pmisc/prometheus/customized"
)

// CollectorRegister model
type CollectorRegister struct {
	Endpoint         string
	JobName          string
	HostName         string
	ETCDURLs         []string
	exit             chan int
	failed_time      int
	collectors       []customized.ICollector
	metricGatewayURL string
	metricAlarmURL   string
	cli              *clientv3.Client
}

var (
	DefaultHostIP string
	errorLogger   *log.Logger
	infoLogger    *log.Logger
	prefix        = "/cruise/"
	CGW           = "/cruise/cgw"
	AGW           = "/cruise/agw"
)

func init() {
	DefaultHostIP = lib.ResolveHostIP()
	errorLogger = log.New(os.Stdout,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

// a func to initial a register of collectors
func NewCollectorRegister(jobname string, ETCDURLs []string) *CollectorRegister {

	// access hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = DefaultHostIP
	}

	// package register param
	cr := &CollectorRegister{
		JobName:  jobname,
		Endpoint: DefaultHostIP,
		ETCDURLs: ETCDURLs,
		HostName: hostname,
	}

	// connect etcd
	cr.connEtcd()
	//watch key
	go func() {
		cr.etcdWatchKey()
	}()

	cr.exit = make(chan int)
	cr.cornTask()

	return cr
}

// registe collector
func (cr *CollectorRegister) Registe(c customized.ICollector) {
	// check whether has same collector
	for _, tc := range cr.collectors {
		if reflect.TypeOf(tc) == reflect.TypeOf(c) {
			return
		}
	}
	cr.collectors = append(cr.collectors, c)
}

func (cr *CollectorRegister) collect() (metrics []customized.Metric) {

	for _, collector := range cr.collectors {
		ms, err, am := collector.Collect()
		if err != nil {
			continue
		}
		// TODO alarm moudle
		if am.MetricName != "" {
			go cr.SendAlarm(am.MetricName, am.Reason, "system")
		}
		metrics = append(metrics, ms...)
	}

	return
}

// Push : push metrics collected by collector to prometheus gateway
func (cr *CollectorRegister) push() error {
	if cr.metricGatewayURL == "" {
		return errors.New("push gateway is empty")
	}
	// package request url http://localhost:9091/metrics/job/%s/instance/%s
	reqURL := fmt.Sprintf("%s/metrics/job/%s/instance/%s/hostname/%s", cr.metricGatewayURL, cr.JobName, cr.Endpoint, cr.HostName)
	metrics := cr.collect()
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
	if cr.failed_time != 0 {
		if cr.failed_time > 60 {
			cr.failed_time = 60
		}
		time.Sleep(5 * time.Duration(cr.failed_time) * time.Second)
	}
	req, err := http.NewRequest("POST", reqURL, bytes.NewBufferString(ms))
	if err != nil {
		cr.failed_time++
		errorLogger.Println(reqURL)

		return err
	}
	resp, err := lib.HttpClient.Do(req)
	if err != nil {
		cr.failed_time++
		errorLogger.Println(reqURL, err.Error())
		return err
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		cr.failed_time++
		return err
	}
	if string(bs) != "" {
		cr.failed_time++
		errorLogger.Println(ms)
		return errors.New(string(bs))
	}
	cr.failed_time = 0
	return nil
}

// it's corn job,collect metrics with a specific duration
func (cr *CollectorRegister) cornTask() {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-ticker.C:
				err := cr.push()
				if err != nil {
					errorLogger.Println("something is wrrong,err reason:", err.Error())
				}
			case <-cr.exit:
				cr.exit = nil
				return
			}
		}
	}()
}

// Exit func ,release occupy resource
func (cr *CollectorRegister) Exit() {
	if cr.exit == nil {
		return
	}
	cr.exit <- 1
}

// to string function
func (cr *CollectorRegister) ToString() string {
	return fmt.Sprintf("Endpoint:%s,JobName:%s,collectors size:%d", cr.Endpoint, cr.JobName, len(cr.collectors))
}

/**
** metricType ,a filed to descript what metric is ,metric is set by user,can be 'cpu','memory' etc....,
** alertType , alertType has two type  ‘system’  and ‘global’  ,default value is 'global'
**** 'system' : hava maximum attemps,be constrained by alarm gateway
**** 'global' : send alarm every time
**/
func (cr *CollectorRegister) SendAlarm(metricType, reason, alertType string) error {
	if cr.metricAlarmURL == "" {
		return errors.New("alarm url is not init")
	}
	// package alarm parameters
	am := customized.AlarmInfo{
		JobName:    cr.JobName,
		HostIP:     cr.Endpoint,
		HostName:   cr.HostName,
		MetricName: metricType,
		Reason:     reason,
		Timestamp:  time.Now().Unix()}

	amBytes, err := json.Marshal(am)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", cr.metricAlarmURL, bytes.NewReader(amBytes))
	if err != nil {
		errorLogger.Fatalln(fmt.Sprintf("occur err when access a http request,url is %s,reason is %s ", cr.metricAlarmURL, err.Error()))
		return err
	}
	resp, err := lib.HttpClient.Do(req)
	if err != nil {
		errorLogger.Fatalln(fmt.Sprintf("occur err when exec client do,url is %s,reason is %s ", cr.metricAlarmURL, err.Error()))
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (cr *CollectorRegister) etcdWatchKey() {
	for {
		if cr.cli == nil {
			cr.connEtcd()
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
	ctx := context.TODO()
	ch := cr.cli.Watch(ctx, prefix, clientv3.WithPrefix())
	for {
		select {
		case c := <-ch:
			for _, e := range c.Events {
				fmt.Println("Watch", string(e.Kv.Key), string(e.Kv.Value))
				if CGW == string(e.Kv.Key) {
					cr.metricGatewayURL = string(e.Kv.Value)
				} else if AGW == string(e.Kv.Key) {
					cr.metricAlarmURL = string(e.Kv.Value)
				}
			}
		}
	}
}

func (cr *CollectorRegister) connEtcd() error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cr.ETCDURLs,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		errorLogger.Println("connection etcd error,reason:", err.Error())
		return err
	}
	infoLogger.Println("connected etcd sucessful")
	ctx := context.TODO()
	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err == nil {
		for i := 0; i < int(resp.Count); i++ {
			infoLogger.Println("GET", string(resp.Kvs[i].Key), string(resp.Kvs[i].Value))
			if CGW == string(resp.Kvs[i].Key) {
				cr.metricGatewayURL = string(resp.Kvs[i].Value)
			} else if AGW == string(resp.Kvs[i].Key) {
				cr.metricAlarmURL = string(resp.Kvs[i].Value)
			}
		}

	} else {
		errorLogger.Println("connection etcd error,reason:", err.Error())
	}
	cr.cli = cli
	return nil
}
