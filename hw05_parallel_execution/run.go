package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	ignoreErrors := m <= 0

	tasksCh := make(chan Task)
	stopCh := make(chan struct{})
	var wg sync.WaitGroup
	var once sync.Once

	var errorCount int64

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case task, ok := <-tasksCh:
					if !ok {
						return
					}

					if err := task(); err != nil && !ignoreErrors {
						newErrorCount := atomic.AddInt64(&errorCount, 1)
						if newErrorCount >= int64(m) {
							once.Do(func() {
								close(stopCh)
							})
							return
						}
					}

				case <-stopCh:
					return
				}
			}
		}()
	}

	go func() {
		defer close(tasksCh)

		for _, task := range tasks {
			select {
			case tasksCh <- task:
			case <-stopCh:
				return
			}
		}
	}()

	wg.Wait()

	if !ignoreErrors && atomic.LoadInt64(&errorCount) >= int64(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
