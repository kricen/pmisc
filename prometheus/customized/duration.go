package customized

type DurationCollector struct {
}

func (c *DurationCollector) Collect() (metrics []Metric, err error) {

	return
}

func (c *DurationCollector) AddRecord(url string, duration int64) {

}
