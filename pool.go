package metric

import (
	"sync"
	"time"
)

type packet struct {
	id     int64 // unique packet identifier
	currNs int64 // current timestamp in nano seconds
	diffNs int64 // timestamp diff in nano seconds
}

func (p *packet) diff() int64 {
	currNs := time.Now().UnixNano()
	return currNs - p.currNs
}

func newPacket(id int64, prevNs int64) packet {
	currNs := time.Now().UnixNano()

	var diffNs int64
	if prevNs != 0 {
		diffNs = currNs - prevNs
	}

	return packet{
		id:     id,
		currNs: currNs,
		diffNs: diffNs,
	}
}

func newPacketCh() chan packet {
	return make(chan packet, 1)
}

type pool struct {
	wg   *sync.WaitGroup
	size int64
}

func NewPool(size int64) pool {
	return pool{
		size: size,
		wg:   new(sync.WaitGroup),
	}
}

func (g *pool) Write() <-chan packet {
	ch := newPacketCh()

	g.wg.Add(1)
	go func() {
		for i := int64(0); i < g.size; i++ {
			ch <- newPacket(i, 0)
		}

		close(ch)
		g.wg.Done()
	}()

	return ch
}

func (g *pool) Read(key string, rch <-chan packet, collect func(b Bucket)) {
	// register read goroutine to save diff results
	bk := newBucket(key, g.size)

	g.wg.Add(1)
	go func() {
		for chPacket := range rch {
			//diff := chPacket.diff()
			bk.Collect(chPacket.id)
			//metrics = append(metrics, diff)
		}

		collect(bk)
		g.wg.Done()
	}()
}

func (g *pool) ReadAndWrite(key string, rch <-chan packet, collect func(b Bucket)) <-chan packet {
	// register read goroutine to save diff results

	wch := newPacketCh()
	bk := newBucket(key, g.size)

	g.wg.Add(1)
	go func() {
		for chPacket := range rch {
			rwPacket := newPacket(chPacket.id, chPacket.currNs)
			bk.Collect(chPacket.id)
			//metrics = append(metrics, rwPacket.diffNs)

			wch <- rwPacket
		}

		collect(bk)
		close(wch)
		g.wg.Done()
	}()

	return wch
}

func (g *pool) Wait() {
	g.wg.Wait()
}
