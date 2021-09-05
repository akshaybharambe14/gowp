package main

import (
	"context"
	"fmt"

	"github.com/akshaybharambe14/gowp"
)

func main() {
	run()
}

func run() {
	const numTasks = 10
	wp, _ := gowp.New(context.TODO(), numTasks, 4, false)

	for i := 0; i < numTasks; i++ {
		i := i
		wp.Submit(func() error {
			fmt.Println("square of ", i, " is ", i*i)
			return nil
		})
	}

	wp.Close()

	_ = wp.Wait()
	fmt.Println("exit")
}
