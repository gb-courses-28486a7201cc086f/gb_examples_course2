package main

import "sync"

type Pool struct {
	mx          sync.Mutex
	wg          *sync.WaitGroup
	jobsQ       chan func()
	jobsQClosed bool
}

// NewPool creates new Pool with {size} workers
func NewPool(size int) *Pool {
	inQ := make(chan func(), size)

	wg := sync.WaitGroup{}
	// create workers
	for i := 0; i < size; i++ {
		wg.Add(1)
		go func(workerId int, inQ chan func()) {
			defer wg.Done()
			for task := range inQ {
				task()
			}
		}(i, inQ)
	}

	return &Pool{wg: &wg, jobsQ: inQ}
}

// Close stops pool and waits until all workers finished
func (p *Pool) Close() {
	p.mx.Lock()
	defer p.mx.Unlock()

	if p.jobsQClosed == false {
		close(p.jobsQ)
		p.jobsQClosed = true
		p.wg.Wait()
	}
}
