package main

import (
	"sync"
	"testing"
)

func newIncrementor() (counterPtr *int, incrementor func()) {
	counter := 0
	incOne := func() {
		counter += 1
	}
	return &counter, incOne
}

// TestIncrementorSingleThread runs += 1 in single workers
// multiply times to check if result is stable
func TestIncrementorSingleThread(t *testing.T) {
	jobsCount := 1000
	expected := 1000

	for i := 0; i < 100; i++ {
		counterPtr, inrementor := newIncrementor()
		pool := NewPool(1)
		// put jobs
		for i := 0; i < jobsCount; i++ {
			pool.jobsQ <- inrementor
		}
		// wait until done
		pool.Close()
		// check
		counter := *counterPtr
		if counter != expected {
			t.Errorf("test %d fail: got %d, expected %d", i, counter, expected)
		}
	}
}

// TestIncrementorMultiThread runs += 1 in some concurent workers
// multiply times to check if result is stable
func TestIncrementorMultiThread(t *testing.T) {
	jobsCount := 1000
	expected := 1000

	for i := 0; i < 100; i++ {
		counterPtr, inrementor := newIncrementor()
		pool := NewPool(jobsCount)
		// put jobs
		for i := 0; i < jobsCount; i++ {
			pool.jobsQ <- inrementor
		}
		// wait until done
		pool.Close()
		// check
		counter := *counterPtr
		if counter != expected {
			t.Errorf("test %d fail: got %d, expected %d", i, counter, expected)
		}
	}
}

func newMutexIncrementor() (counterPtr *int, incrementor func()) {
	mx := sync.Mutex{}
	counter := 0
	incOne := func() {
		mx.Lock()
		defer mx.Unlock()
		counter += 1
	}
	return &counter, incOne
}

// TestIncrementorMultiThread runs += 1 in some concurent workers
// multiply times to check if result is stable
func TestMutexIncrementorMultiThread(t *testing.T) {
	jobsCount := 1000
	expected := 1000

	for i := 0; i < 100; i++ {
		counterPtr, inrementor := newMutexIncrementor()
		pool := NewPool(jobsCount)
		// put jobs
		for i := 0; i < jobsCount; i++ {
			pool.jobsQ <- inrementor
		}
		// wait until done
		pool.Close()
		// check
		counter := *counterPtr
		if counter != expected {
			t.Errorf("test %d fail: got %d, expected %d", i, counter, expected)
		}
	}
}
