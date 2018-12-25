package customized

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/kricen/pmisc/lib"
)

var (
	DefaultHostIP = "127.0.0.1"
	cpuIndex      = -1
	memIndex      = -1
	processID     = strconv.Itoa(os.Getpid())
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

func accesProcessUsage(pid string) (cpuUsage float64, memUsage float64, err error) {

	if runtime.GOOS != "linux" {
		err = errors.New("not support")
		return
	}
	cmd := exec.Command("top", "-b", "-n", "1", "-p", pid)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return
	}
	arr := strings.Split(out.String(), "\n")

	for i := 0; i < len(arr); i++ {
		subArr := strings.Fields(arr[i])
		if strings.Contains(strings.ToUpper(arr[i]), "CPU") && strings.Contains(strings.ToUpper(arr[i]), "MEM") {
			for j := 0; j < len(subArr); j++ {
				if strings.Contains(strings.ToUpper(subArr[j]), "CPU") {
					cpuIndex = j
				}
				if strings.Contains(strings.ToUpper(subArr[j]), "MEM") {
					memIndex = j
				}
			}
		}

		if strings.Contains(arr[i], pid) {
			if memIndex != -1 && cpuIndex != -1 && len(subArr) > memIndex && len(subArr) > cpuIndex {
				cpuUsage, err = strconv.ParseFloat(subArr[cpuIndex], 64)
				if err != nil {
					continue
				}
				memUsage, err = strconv.ParseFloat(subArr[memIndex], 64)
				if err != nil {
					continue
				}

				return
			}
		}
	}

	err = errors.New("notFound")
	return
}
