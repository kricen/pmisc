package customized

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/pmisc/lib"
	"github.com/toolkits/nux"
)

type DiskCollector struct {
	latestAlarmRecord map[string]int64
}

func NewDiskCollector() *DiskCollector {
	record := make(map[string]int64, 0)
	return &DiskCollector{latestAlarmRecord: record}
}

func (c *DiskCollector) Collect() (metrics []Metric, err error, ai AlarmInfo) {
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

		// check metric whether overhead the threshold
		percent := lib.Decimal(du.BlocksUsedPercent / 100)
		if percent > 0.8 {
			lastAlarmTime := c.latestAlarmRecord[du.FsSpec]
			now := time.Now().Unix()
			if now-lastAlarmTime > 5*60 {
				c.latestAlarmRecord[du.FsSpec] = now
				ai.MetricName = "filesystem"
				ai.Reason = fmt.Sprintf("磁盘：%s 使用率超过阈值，现在:%f%s", du.FsSpec, percent*100, "%")
			}
		}
		metrics = append(metrics, Metric{Name: du.FsSpec, Value: lib.Decimal(du.BlocksUsedPercent / 100), NameType: "filesystem", ValueType: reflect.Float64, MetricName: "node_filesystem_used"})
	}

	return
}
