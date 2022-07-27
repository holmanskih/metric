package main

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Latency_Read_Write(t *testing.T) {
	t.Run("write to chan read from chan 1k", func(t *testing.T) {
		testMemReadWrite(t, 1000, "test_r_1k")
	})

	t.Run("write to chan read from chan 10k", func(t *testing.T) {
		testMemReadWrite(t, 10000, "test_r_10k")
	})

	t.Run("write to chan read from chan 100k", func(t *testing.T) {
		testMemReadWrite(t, 100000, "test_r_100k")
	})
}

func Test_Mem_Latency_Read_Double_Write(t *testing.T) {
	t.Run("write to chan read and send to next chan 1k", func(t *testing.T) {
		testMemReadDoubleWrite(t, 1000, "test_d_r1_1k", "test_d_r2_1k")
	})

	t.Run("write to chan read and send to next chan 10k", func(t *testing.T) {
		testMemReadDoubleWrite(t, 10000, "test_d_r1_10k", "test_d_r2_10k")
	})

	t.Run("write to chan read and send to next chan 100k", func(t *testing.T) {
		testMemReadDoubleWrite(t, 100000, "test_d_r1_100k", "test_d_r2_100k")
	})
}

func testMemReadWrite(t *testing.T, n int64, rName1 string) {
	p := NewPool(n)
	m := New(n)

	rch := p.Write()
	p.Read(rName1, rch, func(b Bucket) { m.Collect(b) })
	p.Wait()

	data := m.ActionDiffData()

	var avg int64
	for id, row := range data {
		var sum int64
		for _, cell := range row {
			sum += cell
		}
		avg += sum
		log.Printf("id: %d latency: %d", id, sum)
	}

	avg /= int64(len(data))
	log.Printf("avg latency: %d", avg)

	err := m.ExportToCSV(rName1)
	assert.NoError(t, err)
}

func testMemReadDoubleWrite(t *testing.T, n int64, rName1, rName2 string) {
	p := NewPool(n)
	m := New(n)

	rch := p.Write()
	wch := p.ReadAndWrite(rName1, rch, func(b Bucket) { m.Collect(b) })
	p.Read(rName2, wch, func(b Bucket) { m.Collect(b) })
	p.Wait()

	data := m.ActionDiffData()

	var avg int64
	for id, row := range data {
		var sum int64
		for _, cell := range row {
			sum += cell
		}
		avg += sum
		log.Printf("id: %d latency: %d", id, sum)
	}

	avg /= int64(len(data))
	log.Printf("avg latency: %d", avg)

	err := m.ExportToCSV(rName1 + "+" + rName2)
	assert.NoError(t, err)
}
