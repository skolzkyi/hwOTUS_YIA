package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrErrorsBadErrorsCount = errors.New("errors count less tasks len")
var ErrErrorsBadGoroutinesCount = errors.New("tasks len less goroutines count")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {

	var returnerror error

	if len(tasks) < n {
		returnerror = ErrErrorsBadGoroutinesCount
		return returnerror
	}

	if m > len(tasks) {
		returnerror = ErrErrorsBadErrorsCount
		return returnerror
	}

	workch := make(chan Task, n)
	errch := make(chan error)

	//var workerr error
	wg := sync.WaitGroup{}
	wgStop := sync.WaitGroup{}
	wg.Add(n)
	wgStop.Add(n)
	var stopint int32

	for i := 0; i < n; i++ {
		workch <- tasks[i]
	}

	fmt.Println("count of GR: ", n, "buf: ", len(workch))

	k := 0
	for i := 0; i < n; i++ {
		k++
		go worker(i, workch, errch, &wg, &wgStop, &stopint)

	}
	fmt.Println("real count of GR: ", k)
	k = 0
	q := 0
	t := n
	var flag bool
	for {
		workerr := <-errch
		if workerr != nil {
			if !flag {
				k++
				fmt.Println("numerr: ", k)
			} else {
				q++
				fmt.Println("numerr after stop: ", q)
			}
		}
		if !flag {
			if ((k < m || workerr == nil) || (m <= 0)) && (t <= len(tasks)-1) { //if ((k < m || workerr == nil) && (t <= len(tasks)-1)) || (m <= 0)
				workch <- tasks[t]
				fmt.Println("len: ", t)
				t++
			} else {
				close(workch)
				if (k >= m) && (m > 0) {
					returnerror = ErrErrorsLimitExceeded
				}
				flag = true
			}
		}
		fmt.Println("stopint: ", stopint)
		if stopint >= int32(n) {
			fmt.Println("break")
			break
		}
	}
	fmt.Println("freedom")
	wg.Wait()

	return returnerror
}

func worker(i int, workch chan Task, errch chan error, wg *sync.WaitGroup, wgStop *sync.WaitGroup, stopint *int32) {
	//j := 1
	defer wg.Done()
	//fmt.Println("start worker: ", i)
	for {
		task, ok := <-workch
		if !ok {
			//fmt.Println("break worker: ", i)
			atomic.AddInt32(stopint, int32(1))
			//wgStop.Done()
			break
		} else {
			//fmt.Println("worker: ", i, " task: ", j, "buf: ", len(workch))
			//j++
			err := task()
			errch <- err
		}
	}
	//fmt.Println("stop worker: ", i)
}
