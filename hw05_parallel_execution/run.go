package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {

	workch := make(chan Task, len(tasks))
	errch := make(chan error)
	signalch := make(chan struct{})
	wg := sync.WaitGroup{}

	for i := 0; i <= len(tasks); i++ {
		workch <- tasks[i]
	}
	close(workch)

	wg.Add(n + 1)

	go watchdog(n, workch, errch, signalch, wg)(n, workch, errch, signalch, wg)
	<-signalch

	for i := 0; i < n; i++ {
		go worker(tasks, workch, errch, wg)(tasks, workch, errch, wg)
	}

	wg.Wait()

	return nil
}

func worker(tasks []Task, workch chan Task, errch chan error, wg *sync.WaitGroup) {
	for {
		task, ok := <-workch
		err := task()
		if err != nil {
			errch <- err
		}
		if !ok {
			break
		}
		wg.Done()
	}

}

func watchdog(n int, workch chan Task, errch chan error, signalch chan struct{}, wg *sync.WaitGroup) {

	signalch <- struct{}{}

	i := 0
	for i < n {
		<-errch
		i++
	}
	for {
		_, ok := <-workch
		if !ok {
			break
		}

	}
	wg.Done()
}
