package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
)

func isPrime(value int) bool {
	// well known values
	if value <= 1 {
		return false
	}
	if value == 2 {
		return true
	}
	// Eratosphene method used
	// make [0,0,2,3,4,...,value] sequence...
	// Eratosphene needs [2,3,4,...,value],
	// but we add leading zeros to simplify: index == value
	primeCandidates := make([]int, value+1)
	for i := 2; i <= value; i++ {
		primeCandidates[i] = i
	}
	for step := 2; step*step <= value; {
		// numbers less step^2 checked on prev iteration
		// so, skip it
		for i := step * step; i <= value; i += step {
			if i%step == 0 {
				primeCandidates[i] = 0
			}
		}
		// if {value} was dropped - it is not prime
		if primeCandidates[value] == 0 {
			return false
		}
		// new step is first non "dropped" more thn current step
		for i := step + 1; i < value; i++ {
			if primeCandidates[i] != 0 {
				step = i
				break
			}
		}
	}
	// {value} was not dropped => it is prime
	return primeCandidates[value] != 0
}

func checkPrimes(first, last int) {
	for i := first; i < last; i++ {
		fmt.Printf("%d is prime? %v\n", i, isPrime(i))
		// checking "is prime" it is CPU-bound
		// so, we use Gosched() to allow all workers to run
		runtime.Gosched()
	}
}

func main() {
	// tracing
	trace.Start(os.Stderr)
	defer trace.Stop()

	// demo: try to check some numbers using 4 workers on 2 threads
	runtime.GOMAXPROCS(2)
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		defer wg.Done()
		checkPrimes(1, 25)
	}()
	go func() {
		defer wg.Done()
		checkPrimes(26, 50)
	}()
	go func() {
		defer wg.Done()
		checkPrimes(51, 75)
	}()
	go func() {
		defer wg.Done()
		checkPrimes(76, 100)
	}()
	wg.Wait()
}
