package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrErrorsBadGoroutinesCount = errors.New("goroutines count large tasks len")
var ErrErrorsBadErrorsCount = errors.New("errors count large tasks len")
var ChannelCounter int

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {

	ChannelCounter = 0
	var returnerror error

	if n > len(tasks) {
		returnerror = ErrErrorsBadGoroutinesCount
		return returnerror
	}

	if m > len(tasks) {
		returnerror = ErrErrorsBadErrorsCount
		return returnerror
	}

	workch := make(chan Task, n)
	errch := make(chan error, n)

	var stopint int32
	var rwm sync.RWMutex

	for i := 0; i < n; i++ {
		workch <- tasks[i]
	}

	//	fmt.Println("count of GR: ", n, "buf: ", len(workch))

	k := 0
	for i := 0; i < n; i++ {
		k++
		go worker(i, n, workch, errch, &rwm, &stopint)

	}
	//	fmt.Println("real count of GR: ", k)
	k = 0
	q := 0
	t := n
	var flag bool

	for {

		workerr := <-errch
		if workerr != nil {
			if !flag {
				k++
				//	fmt.Println("numerr: ", k)
			} else {
				q++
				//	fmt.Println("numerr after stop: ", q)
			}
		}

		if !flag {
			if ((k < m || workerr == nil) || (m <= 0)) && (t < len(tasks)) { //if ((k < m || workerr == nil) && (t <= len(tasks)-1)) || (m <= 0)
				workch <- tasks[t]
				//fmt.Println("len: ", t)
				t++
			} else {
				//fmt.Println("workch closed!")
				close(workch)
				if (k >= m) && (m > 0) {
					returnerror = ErrErrorsLimitExceeded
				}
				flag = true
			}
		}

		rwm.RLock()
		//	fmt.Println("ChannelCounter: ", ChannelCounter, " n ", n)
		if ChannelCounter >= n {
			rwm.RUnlock()
			//fmt.Println("break")
			break
		}
		rwm.RUnlock()

	}
	//fmt.Println("freedom")

	return returnerror
}

//cd C:\REPO\Go\!OTUS\hwOTUS_YIA\hw05_parallel_execution

func worker(i int, n int, workch chan Task, errch chan error, rwm *sync.RWMutex, stopint *int32) {

	for {
		task, ok := <-workch
		if !ok {
			//fmt.Println("break worker: ", i)

			rwm.Lock()
			if ChannelCounter == n-1 {
				close(errch)
			}
			ChannelCounter++
			//	fmt.Println("ChannelCounter: ", ChannelCounter, " from worker ", i)
			rwm.Unlock()

			break
		} else {

			err := task()
			errch <- err
		}
	}

}
