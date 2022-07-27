package metric

import (
	"time"
)

type Bucket interface {
	// Name returns unique identifier bucket name
	Name() string

	// Collect collects timestamp data and assign it to id
	Collect(id int64)

	// Metric returns all bucket`s metric data in form of timestamps
	Metric() []int64
}

type bucket struct {
	name   string
	size   int64
	metric []int64
}

func (b *bucket) Name() string {
	return b.name
}

func (b *bucket) Metric() []int64 {
	return b.metric
}

func (b *bucket) Collect(id int64) {
	if id >= b.size {
		return
	}

	ts := time.Now().UnixNano()
	b.metric[id] = ts
}

func newBucket(key string, size int64) Bucket {
	return &bucket{
		name:   key,
		size:   size,
		metric: make([]int64, size),
	}
}
