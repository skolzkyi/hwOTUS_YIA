package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded      = errors.New("errors limit exceeded")
	ErrErrorsBadGoroutinesCount = errors.New("goroutines count large tasks len")
	ErrErrorsBadErrorsCount     = errors.New("errors count large tasks len")
	ChannelCounter              int
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	ChannelCounter = 0
	var returnerror error

	returnerror = checkinputParameters(len(tasks), n, m)
	if returnerror != nil {
		return returnerror
	}

	workch := make(chan Task, n)
	errch := make(chan error, n)

	var rwm sync.RWMutex

	for i := 0; i < n; i++ {
		workch <- tasks[i]
	}

	k := 0
	for i := 0; i < n; i++ {
		k++
		go worker(n, workch, errch, &rwm)
	}

	k = 0
	t := n
	var flag bool

	for {
		workerr := <-errch
		if workerr != nil {
			if !flag {
				k++
			}
		}
		if !flag {
			if ((k < m || workerr == nil) || (m <= 0)) && (t < len(tasks)) {
				workch <- tasks[t]
				t++
			} else {
				close(workch)
				if (k >= m) && (m > 0) {
					returnerror = ErrErrorsLimitExceeded
				}
				flag = true
			}
		}

		rwm.RLock()
		if ChannelCounter >= n {
			rwm.RUnlock()
			break
		}
		rwm.RUnlock()
	}

	return returnerror
}

func worker(n int, workch chan Task, errch chan error, rwm *sync.RWMutex) {
	for {
		task, ok := <-workch
		if !ok {
			rwm.Lock()
			if ChannelCounter == n-1 {
				close(errch)
			}
			ChannelCounter++
			rwm.Unlock()
			break
		} else {
			err := task()
			errch <- err
		}
	}
}

func checkinputParameters(l, n, m int) error {
	var returnerror error
	if n > l {
		returnerror = ErrErrorsBadGoroutinesCount
		return returnerror
	}

	if m > l {
		returnerror = ErrErrorsBadErrorsCount
		return returnerror
	}

	return nil
}
