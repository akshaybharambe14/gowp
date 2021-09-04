/*
Package gowp provides a pool of workers with limited concurrency.

Example:
	// create a pool for 10 tasks with 4 workers that is exists on error
	wp := gowp.New(context.TODO(), 10, 4, true)

	// add tasks to the pool
	for i := 0; i < 10; i++ {
		wp.Submit(func() error {
			// do something
			if i%2 == 0 {
				// return error if something happens
				return errors.New("error")
			}

			return nil
		}()
	}

	// signal that we are done producing work
	wp.Close()

	// wait for all the workers to finish their work, check the first error, if any
	if err := wp.Wait(); err != nil {
		// handle error
	}
*/
package gowp
