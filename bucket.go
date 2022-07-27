package main

import (
	"log"
	"time"
)

type Bucket interface {
	Name() string
	Add(id int64)
	Metric() []int64
}

type bucket struct {
	key    string
	size   int64
	metric []int64
}

func (b *bucket) Name() string {
	return b.key
}

func (b *bucket) Metric() []int64 {
	return b.metric
}

func (b *bucket) Add(id int64) {
	if id >= b.size {
		return
	}

	ts := time.Now().UnixNano()
	log.Printf("bucket: %s id: %d ts: %d", b.key, id, ts)
	b.metric[id] = ts
}

func newBucket(key string, size int64) Bucket {
	return &bucket{
		key:    key,
		size:   size,
		metric: make([]int64, size),
	}
}
