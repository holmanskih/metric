package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Collector interface {
	Size() int64
	Collect(b Bucket)
	ActionDiffData() [][]int64 // returns metric data with id and ts diff in nanoseconds
	PrevActionDiffData() [][]int64
	ExportToCSV(name string) error
}

type collector struct {
	size    int64    // bucket size
	buckets []Bucket // metric buckets
}

func (c *collector) Size() int64 {
	return c.size
}

func (c *collector) Collect(b Bucket) {
	c.buckets = append(c.buckets, b)
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

func (c *collector) ExportToCSV(name string) error {
	actionDiff := c.ActionDiffData()
	if err := c.exportCSV(name+"_actionDiff", actionDiff); err != nil {
		return err
	}

	prevaActionDiff := c.PrevActionDiffData()
	if err := c.exportCSV(name+"_prevActionDiff", prevaActionDiff); err != nil {
		return err
	}
	return nil
}

func (c *collector) exportCSV(name string, metric [][]int64) error {
	path := fmt.Sprintf("%s.csv", name)

	cols := len(c.buckets)
	headers := make([]string, 0, cols)
	for _, b := range c.buckets {
		headers = append(headers, b.Name())
	}

	if err := newCsv(path, headers, metric, cols); err != nil {
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
		csvData := make([]string, len(metricData)+1)
		csvData[0] = strconv.FormatInt(int64(id), 10)

		var sum int64
		for i, data := range metricData {
			csvData[i+1] = strconv.FormatInt(data, 10)
			sum += data
		}

		//log.Printf("avg case: %s value: %d", id, sum/int64(len(metricData)))

		if err = writer.Write(csvData); err != nil {
			return err
		}
	}

	return nil
}

func New(size int64) Collector {
	return &collector{
		size:    size,
		buckets: make([]Bucket, 0),
	}
}
