package customized

import (
	"log"
	"reflect"

	"github.com/pmisc/lib"
	"github.com/toolkits/nux"
)

type DiskCollector struct {
}

func NewDiskCollector() *DiskCollector {
	return &DiskCollector{}
}

func (c *DiskCollector) Collect() (metrics []Metric, err error) {
	mountPoints, err := nux.ListMountPoint()
	if err != nil {
		log.Printf("collecting disk info err, reason:%s", err.Error())
		return
	}

	for idx := range mountPoints {
		var du *nux.DeviceUsage
		du, err = nux.BuildDeviceUsage(mountPoints[idx][0], mountPoints[idx][1], mountPoints[idx][2])
		if du.FsSpec == "proc" {
			continue
		}
		if err != nil {
			continue
		}
		metrics = append(metrics, Metric{Name: du.FsSpec, Value: lib.Decimal(du.BlocksUsedPercent / 100), NameType: "filesystem", ValueType: reflect.Float64, MetricName: "node_filesystem_used"})
	}

	return
}
