package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrInvalidArgs         = errors.New("one goroutine (N) and one possible error (M) required to complete tasks")
)

type Task func() error

func worker(tasks <-chan Task, done <-chan struct{}, remainErrCount *int64, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-done:
			return
		default:
			if task, ok := <-tasks; ok && atomic.LoadInt64(remainErrCount) > 0 {
				if task() != nil {
					// уменьшаем счетчик допустимых ошибок
					atomic.AddInt64(remainErrCount, -1)
				}
			} else {
				return
			}
		}
	}
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, N, M int) error {
	if N <= 0 || M <= 0 {
		return ErrInvalidArgs
	}

	if len(tasks) < N {
		N = len(tasks)
	}

	donech := make(chan struct{})
	tasksch := make(chan Task, len(tasks))

	remainErrCount := int64(M)

	wg := &sync.WaitGroup{}
	wg.Add(N)
	for i := 0; i < N; i++ {
		go worker(tasksch, donech, &remainErrCount, wg)
	}

	for _, task := range tasks {
		// добавляем в очередь, если не превышено число ошибок
		if atomic.LoadInt64(&remainErrCount) > 0 {
			tasksch <- task
			continue
		}

		donech <- struct{}{}
		return ErrErrorsLimitExceeded
	}

	close(tasksch)
	wg.Wait()
	close(donech)

	if atomic.LoadInt64(&remainErrCount) <= 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
