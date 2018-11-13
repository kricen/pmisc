package registry

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"

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
			metrics, err := rc.Collect()
			if err != nil {
				fmt.Printf("error:%s", err.Error())
			}
			for _, metric := range metrics {
				fmt.Println(metric.Name, metric.Value, metric.NameType, metric.ValueType, metric.MetricName)
			}
		}
	}()
	wg.Wait()
	fmt.Println("--------------", "--------------")
	metrics, err := rc.Collect()
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

	cr := NewCollectorRegister("monitor-helper", "url")

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
		cr.Push()
		metrics := cr.Collect()
		fmt.Printf("-------------info:%d---------------\n", i)
		fmt.Printf("Endpoint:%s ,JobName:%s\n", cr.Endpoint, cr.JobName)
		for _, metric := range metrics {
			fmt.Println(metric.Name, metric.Value, metric.NameType, metric.ValueType, metric.MetricName)
		}
	}

}

type Hello struct {
	Name string
}

func (h Hello) SetName(name string) {
	h.Name = name
	fmt.Println("hello:", &h, name)
	fmt.Printf("%p", &h)
}
func (h Hello) GetName() string {

	return h.Name
}

func TestStruct(t *testing.T) {

	hello := &Hello{}
	fmt.Printf("%p", &hello)
	hello.SetName("hood")
	fmt.Printf("--------,%s,---\n", hello.GetName())
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		hello.SetName("good")
		time.Sleep(3 * time.Second)
		fmt.Println(hello.GetName())

	}()
	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Second)

		hello.SetName("good2")

	}()
	wg.Wait()

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

func TestUnix(t *testing.T) {
	startAt := time.Now().UnixNano() / 1000 / 1000
	time.Sleep(1 * time.Second)
	endAt := time.Now().UnixNano() / 1000 / 1000
	t.Log(endAt - startAt)

}
