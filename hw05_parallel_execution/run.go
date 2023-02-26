package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrBadWorkersCount     = errors.New("workers count less one")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var err error
	var errCounter int
	var wg sync.WaitGroup
	var rwm sync.RWMutex

	if n < 1 {
		err = ErrBadWorkersCount
		return err
	}

	workch := make(chan Task)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(workch chan Task, wg *sync.WaitGroup, rwm *sync.RWMutex) {
			for {
				task, ok := <-workch
				if !ok {
					break
				}
				err := task()
				if err != nil {
					rwm.Lock()
					errCounter++
					rwm.Unlock()
				}
			}
			wg.Done()
		}(workch, &wg, &rwm)
	}

	for _, task := range tasks {
		workch <- task
		rwm.RLock()
		if errCounter >= m && m > 0 {
			rwm.RUnlock()
			close(workch)
			err = ErrErrorsLimitExceeded
			break
		}
		rwm.RUnlock()
	}
	if err == nil {
		close(workch)
	}

	wg.Wait()

	return err
}
