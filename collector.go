package metric

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Collector interface {
	// NewBucket creates new metric bucket. Should be created separately for working inside goroutines
	NewBucket(name string) Bucket

	// Size is max size of metric data to be collected
	Size() int64

	// Collect add buckets metric data to collector
	Collect(b ...Bucket)

	// ActionDiffData returns metric data with id and ts diff in nanoseconds
	ActionDiffData() [][]int64
	PrevActionDiffData() [][]int64

	// ExportToCSV exports collected bucket data to separate folder with .csv files
	ExportToCSV() error
}

const (
	exportBasePath = "_metric"
)

type collector struct {
	name    string   // collector name, uses during metric exporting
	size    int64    // bucket size
	buckets []Bucket // metric buckets
}

func (c *collector) NewBucket(name string) Bucket {
	return newBucket(name, c.size)
}

func (c *collector) Size() int64 {
	return c.size
}

func (c *collector) Collect(b ...Bucket) {
	c.buckets = append(c.buckets, b...)
}

func (c *collector) PrevActionDiffData() [][]int64 {
	data := make([][]int64, c.size)

	for i := int64(0); i < c.size; i++ {

		diffs := make([]int64, 0)
		for _, b := range c.buckets {
			// check first bucket
			if i == 0 {
				diffs = append(diffs, 0)
				continue
			}

			prev := b.Metric()[i-1]
			curr := b.Metric()[i]
			diff := curr - prev
			diffs = append(diffs, diff)
		}

		data[i] = diffs
	}

	return data
}

func (c *collector) ActionDiffData() [][]int64 {
	data := make([][]int64, c.size)

	for i := int64(0); i < c.size; i++ {

		diffs := make([]int64, 0)
		for j := range c.buckets {
			// check first bucket
			if j == 0 {
				diffs = append(diffs, 0)
				continue
			}

			prev := c.buckets[j-1].Metric()[i]
			curr := c.buckets[j].Metric()[i]
			diff := curr - prev
			diffs = append(diffs, diff)
		}

		data[i] = diffs
	}

	return data
}

func (c *collector) ExportToCSV() error {
	ts := strconv.FormatInt(time.Now().Unix(), 10)

	if _, err := os.Stat(exportBasePath); os.IsNotExist(err) {
		if err = os.Mkdir(exportBasePath, 0755); err != nil {
			return err
		}
	}

	path := c.name + "_" + ts
	path = filepath.Join(exportBasePath, path)
	if err := os.Mkdir(path, 0755); err != nil {
		return err
	}

	actionDiff := c.ActionDiffData()
	prevActionDiff := c.PrevActionDiffData()

	file := filepath.Join(path, "action_diff")
	if err := c.exportCSV(file, actionDiff); err != nil {
		return err
	}
	file = filepath.Join(path, "prev_action_diff")
	if err := c.exportCSV(file, prevActionDiff); err != nil {
		return err
	}

	return nil
}

func (c *collector) exportCSV(name string, metric [][]int64) error {
	cols := len(c.buckets)
	headers := make([]string, 0, cols)
	for _, b := range c.buckets {
		headers = append(headers, b.Name())
	}

	if err := newCsv(name+".csv", headers, metric, cols); err != nil {
		return err
	}

	return nil
}

func newCsv(path string, headers []string, metric [][]int64, n int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// headers
	csvData := make([]string, n+1)
	csvData[0] = "id"

	for i := 0; i < n; i++ {
		csvData[i+1] = headers[i]
	}

	if err = writer.Write(csvData); err != nil {
		return err
	}

	// data
	for id, metricData := range metric {
		rowData := make([]string, len(metricData)+1)
		rowData[0] = strconv.FormatInt(int64(id), 10)

		var sum int64
		for i, data := range metricData {
			rowData[i+1] = strconv.FormatInt(data, 10)
			sum += data
		}

		if err = writer.Write(rowData); err != nil {
			return err
		}
	}

	return nil
}

func Init(size int64, name string) Collector {
	return &collector{
		size:    size,
		name:    name,
		buckets: make([]Bucket, 0),
	}
}
