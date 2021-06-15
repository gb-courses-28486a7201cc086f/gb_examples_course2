package workerpool

import (
	"errors"
	"log"
	"testing"
	"time"
)

const (
	poolSize  = 4
	jobsCount = poolSize * 2
)

type TestJob struct {
	id    int
	delay time.Duration
	done  bool
}

func (tj *TestJob) Run() {
	time.Sleep(tj.delay)
	tj.done = true
}

func Example() {
	// create pool
	workers := 10
	pool, err := NewPool(workers)
	if err != nil {
		log.Println(err)
	}

	// make jobs
	jobs := make([]Job, 0, jobsCount)
	for i := 0; i < jobsCount; i++ {
		delay := 5 * time.Millisecond
		jobs = append(jobs, &TestJob{i, delay, false})
	}

	// run jobs
	err = pool.RunBatch(jobs...)
	if err != nil {
		log.Println(err)
	}

	// take results
	for _, job := range jobs {
		job := job.(*TestJob)
		log.Printf("job id=%d done: %v", job.id, job.done)
	}

	// stop pool
	pool.Join()
}

func poolBatchTest(t *testing.T, pool Pool) {
	jobs := make([]Job, 0, jobsCount)
	for i := 0; i < jobsCount; i++ {
		delay := 5 * time.Millisecond
		jobs = append(jobs, &TestJob{i, delay, false})
	}

	err := pool.RunBatch(jobs...)
	if err != nil {
		t.Errorf("run batch error: %v", err)
	}

	for _, job := range jobs {
		testJob := job.(*TestJob)
		if testJob.done != true {
			t.Errorf("job id=%d did not complete", testJob.id)
		}
	}
}

func TestPool(t *testing.T) {
	// setup
	pool, err := NewPool(poolSize)
	if err != nil {
		t.Fatalf("pool creation fail: %v", err)
	}

	// teardown
	defer pool.Join()

	// tests
	t.Run("first batch", func(t *testing.T) {
		poolBatchTest(t, *pool)
	})
	t.Run("second batch", func(t *testing.T) {
		poolBatchTest(t, *pool)
	})
	t.Run("batch after join", func(t *testing.T) {
		jobs := make([]Job, 0, jobsCount)
		for i := 0; i < jobsCount; i++ {
			delay := 5 * time.Millisecond
			jobs = append(jobs, &TestJob{i, delay, false})
		}

		pool.Join()
		err := pool.RunBatch(jobs...)
		if !errors.Is(err, PoolClosedErr) {
			t.Errorf("got %v error, expected: %v", err, PoolClosedErr)
		}
	})
}

func TestPoolBadSize(t *testing.T) {
	_, err := NewPool(-1)
	if !errors.Is(err, PoolNegativeSizeErr) {
		t.Errorf("got %v error, expected: %v", err, PoolNegativeSizeErr)
	}
}
