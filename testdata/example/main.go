package main

import (
	"context"
	"sync"
	"time"

	"github.com/holmanskih/metric"
)

var (
	col metric.Collector
	in  metric.Bucket
	out metric.Bucket
)

func init() {
	//  init metric collector first to get size value first that is used in bucket initialization
	col = metric.Init(50, "example")

	in = col.NewBucket("in")
	out = col.NewBucket("out")
}

func main() {
	i := int64(0)
	wg := sync.WaitGroup{}
	idCh := make(chan int64)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	wg.Add(1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(idCh)
				wg.Done()
				return

			default:
				in.Collect(i)
				idCh <- i
				i++
			}
		}
	}()

	for id := range idCh {
		out.Collect(id)
	}

	wg.Wait()
	col.Collect(in, out)
	_ = col.ExportToCSV()
}
