package benchmarks

import (
	"testing"

	"github.com/akshaybharambe14/gowp"
	"github.com/alitto/pond"
	"github.com/gammazero/workerpool"
)

var (
	noOpErr = func() error {
		return nil
	}

	noOp = func() {
	}
)

const (
	numTasks10  = 10
	numWorkers4 = 4
)

func simple_gowp() {
	wp, _ := gowp.New(numTasks10, gowp.WithNumWorkers(numWorkers4))

	for i := 0; i < numTasks10; i++ {
		wp.Submit(noOpErr)
	}

	// wp.Close()

	_ = wp.Wait()
}

func simple_workerpool() {
	wp := workerpool.New(numWorkers4)

	for i := 0; i < numTasks10; i++ {
		wp.Submit(noOp)
	}

	wp.StopWait()
}

func simple_pond() {
	pool := pond.New(numWorkers4, numTasks10)

	for i := 0; i < numTasks10; i++ {
		pool.Submit(noOp)
	}

	pool.StopAndWait()
}

func Benchmark_simple_gowp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simple_gowp()
	}
}

func Benchmark_simple_workerpool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simple_workerpool()
	}
}

func Benchmark_simple_pond(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simple_pond()
	}
}
