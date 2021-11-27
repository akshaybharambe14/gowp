package gowp_test

import (
	"errors"
	"fmt"
	"time"

	"github.com/akshaybharambe14/gowp"
)

// ExitOnError demonstrates the use of a Worker Pool to implement a use-case
// where a task can fail and the pool should exit and return the error.
func ExamplePool_exitOnError() {
	const (
		numJobs    = 50
		numWorkers = 2
		closeOnErr = true
	)

	wp, _ := gowp.New(
		numJobs,
		gowp.WithExitOnError(true),
		gowp.WithNumWorkers(numWorkers),
	)

	for i := 0; i < numJobs; i++ {
		i := i
		_ = wp.Submit(func() error {
			fmt.Println("processing ", i)
			time.Sleep(time.Millisecond)
			if i == 2 {
				return errors.New("can't continue")
			}
			return nil
		})
	}

	if err := wp.Wait(); err != nil {
		fmt.Println("process jobs: ", err)
	}

	// Unordered output: 50
	// processing 0
	// processing 2
	// processing 1
	// ...
	// process jobs: can't continue
}
