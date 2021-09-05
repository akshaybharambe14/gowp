package benchmarks

import (
	"context"
	"testing"

	"github.com/akshaybharambe14/gowp"
	"github.com/alitto/pond"
	"github.com/gammazero/workerpool"
)

var (
	noOpGoWP = func() error {
		return nil
	}

	noOpWorkerpool = func() {
	}
)

func GoWPSimple() {
	const numTasks = 10
	wp, _ := gowp.New(context.TODO(), numTasks, 4, false)

	for i := 0; i < numTasks; i++ {
		wp.Submit(noOpGoWP)
	}

	wp.Close()

	_ = wp.Wait()
}

func WorkerpoolSimple() {
	const numTasks = 10
	wp := workerpool.New(4)

	for i := 0; i < numTasks; i++ {
		wp.Submit(noOpWorkerpool)
	}

	wp.StopWait()
}

func PondSimple() {
	const numTasks = 10
	pool := pond.New(4, numTasks)

	for i := 0; i < numTasks; i++ {
		pool.Submit(noOpWorkerpool)
	}

	pool.StopAndWait()
}

func BenchmarkOwn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GoWPSimple()
	}
}

func BenchmarkWorkerPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WorkerpoolSimple()
	}
}

func BenchmarkPond(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PondSimple()
	}
}
