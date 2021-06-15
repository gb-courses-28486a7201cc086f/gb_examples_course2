// Package implements tools for running any set of jobs
// in worker pool.
//
// For creating new pool use NewPool(size) method.
// Implement you custom jobs with .Run() method
// and put it into pool using pool.RunBatch(jobs...)
package workerpool

import (
	"errors"
	"sync"
)

var (
	PoolClosedErr       = errors.New("pool closed")
	PoolNegativeSizeErr = errors.New("pool size cannot be negative")
)

// Job is a single unit of work.
// Job should be self-contained
type Job interface {
	Run()
}

type workerTask struct {
	payload Job
	done    chan<- struct{}
}

// worker is a goroutine which takes
// new tasks from queue and process it
type worker struct {
	id    int
	wg    *sync.WaitGroup
	tasks <-chan workerTask
}

func (w *worker) handle() {
	defer w.wg.Done()
	for task := range w.tasks {
		task.payload.Run()
		task.done <- struct{}{}
	}
}

// Pool is manager of set of workers (goroutines),
// it handle input queue for tasks
type Pool struct {
	size    int
	wg      *sync.WaitGroup
	mx      *sync.RWMutex
	inQueue chan workerTask
	closed  bool
}

// Join method used for gracefully pool shutdown.
// It returns after all workers finished work
func (p *Pool) Join() {
	p.mx.Lock()
	defer p.mx.Unlock()
	// prevent panic on close twice
	if !p.closed {
		p.closed = true
		close(p.inQueue)
		p.wg.Wait()
	}
}

// RunBatch put set of jobs to workers queue
// and returns after all those jobs will be done
func (p *Pool) RunBatch(jobs ...Job) error {
	p.mx.RLock()
	if p.closed {
		p.mx.RUnlock()
		return PoolClosedErr
	}
	doneChan := make(chan struct{}, len(jobs))
	// add tasks...
	for _, job := range jobs {
		p.inQueue <- workerTask{
			payload: job,
			done:    doneChan,
		}
	}
	// release lock - we can close chan
	// when tasks has been send
	p.mx.RUnlock()
	// ...and wait until all done
	for i := 0; i < len(jobs); i++ {
		<-doneChan
	}
	return nil
}

// NewPool creates new pool with defined count of workers
func NewPool(size int) (*Pool, error) {
	if size <= 0 {
		return nil, PoolNegativeSizeErr
	}

	wg := &sync.WaitGroup{}
	mx := &sync.RWMutex{}
	inputQueue := make(chan workerTask, size)

	for i := 0; i < size; i++ {
		worker := &worker{i, wg, inputQueue}
		wg.Add(1)
		go worker.handle()
	}

	return &Pool{size: size, wg: wg, mx: mx, inQueue: inputQueue}, nil
}
