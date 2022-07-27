package main

var (
	Counter int64
	col     Collector
)

func MetricCollector(size int64) Collector {
	if col == nil {
		col = New(size)
		return col
	}

	return col
}

func NewBucket(name string) Bucket {
	if col == nil {
		panic("metric collector was not initialized")
	}

	return newBucket(name, col.Size())
}
