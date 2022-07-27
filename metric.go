package metric

var (
	ID  int64
	col Collector
)

// Init initialize Collector with max size of metric data to collect and name for exporting data
func Init(size int64, name string) Collector {
	if col == nil {
		col = newCollector(size, name)
		return col
	}

	return col
}

// NewBucket creates new metric bucket. Should be created separately for working inside goroutines
func NewBucket(name string) Bucket {
	if col == nil {
		panic("metric collector was not initialized")
	}

	return newBucket(name, col.Size())
}
