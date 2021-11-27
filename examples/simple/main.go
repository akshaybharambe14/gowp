package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/akshaybharambe14/gowp"
)

func main() {
	const (
		numJobs    = 50
		numWorkers = 2
		closeOnErr = true
	)

	wp, _ := gowp.New(numJobs, gowp.WithExitOnError(true), gowp.WithNumWorkers(numWorkers))

	for i := 0; i < numJobs; i++ {
		i := i
		_ = wp.Submit(func() error {
			fmt.Println("processing", i)
			time.Sleep(time.Millisecond)
			if i == 12 {
				return errors.New("can't continue")
			}
			return nil
		})
	}

	if err := wp.Wait(); err != nil {
		fmt.Println("process jobs: ", err)
	}
}
