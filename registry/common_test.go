package registry

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/pmisc/lib"
	"github.com/pmisc/prometheus/customized"
)

func TestCollector(t *testing.T) {
	var wg sync.WaitGroup
	rc := &customized.RequestCollector{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			rc.AddRecord("hello", int64(i))
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 2; i++ {
			time.Sleep(2 * time.Second)
			metrics, err, _ := rc.Collect()
			if err != nil {
				fmt.Printf("error:%s", err.Error())
			}
			for _, metric := range metrics {
				fmt.Println(metric.Name, metric.Value, metric.NameType, metric.ValueType, metric.MetricName)
			}
		}
	}()
	wg.Wait()
	metrics, err, _ := rc.Collect()
	if err != nil {
		fmt.Printf("error:%s", err.Error())
	}
	for _, metric := range metrics {
		fmt.Println(metric.Name, metric.Value, metric.NameType, metric.ValueType, metric.MetricName)
	}

}

func TestRegistry(t *testing.T) {

	cc := customized.NewCpuCollector()
	dc := customized.NewDiskCollector()
	mc := customized.NewMemoryCollector()

	rc := customized.NewRequestCollector()

	cr, _ := NewCollectorRegister("monitor-helper", []string{"url"})

	cr.Registe(cc)
	cr.Registe(dc)
	cr.Registe(mc)
	cr.Registe(rc)

	fmt.Println(cr.ToString())
	go func() {
		for i := 0; i < 100; i++ {
			time.Sleep(100 * time.Millisecond)
			rc.AddRecord("getUserName", int64(i))
		}
	}()

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		cr.push()
		metrics := cr.collect()
		fmt.Printf("-------------info:%d---------------\n", i)
		fmt.Printf("Endpoint:%s ,JobName:%s\n", cr.Endpoint, cr.JobName)
		for _, metric := range metrics {
			fmt.Println(metric.Name, metric.Value, metric.NameType, metric.ValueType, metric.MetricName)
		}
	}

}

func TestPost(t *testing.T) {
	var text string
	// for i := 0; i < 10; i++ {
	// 	text += fmt.Sprintf("hello%d_metric %d.14\n", i, i)
	// }
	text = `hellos{method="post",quantile="0.9"} 0.000104474`

	fmt.Println(text)
	req, err := http.NewRequest("POST", "http://47.104.62.159:9091/metrics/job/hello", bytes.NewBufferString(text))
	if err != nil {
		fmt.Println("init req err:", err.Error())
		return
	}
	resp, err := lib.HttpClient.Do(req)
	if err != nil {
		fmt.Println("do req err:", err.Error())
		return
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(bs))
	resp.Body.Close()

}

func TestEtcd(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		t.Log(err.Error())
		return
	}

	ctx := context.TODO()

	resp, err := cli.Get(ctx, "/key/", clientv3.WithPrefix())
	if err == nil {
		fmt.Println(resp.Count)
		for i := 0; i < int(resp.Count); i++ {
			fmt.Println(string(resp.Kvs[i].Key), string(resp.Kvs[i].Value))
		}
	}
	ch := cli.Watch(ctx, "/key/", clientv3.WithPrefix())
	for {
		log.Print("rev")
		select {
		case c := <-ch:
			for _, e := range c.Events {
				log.Printf("%+v", e)
			}
		}
	}

}
