package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/trace"
	"sync"
	"time"

	"geekbrains/examples/lesson6/hw/set"
)

const testSliceSize = 1_000_000

// fills two slices by some random values
func createTestSlices() [][]int {
	rand.Seed(time.Now().Unix())
	origin := make([]int, testSliceSize)
	for i := 0; i < testSliceSize; i++ {
		origin[i] = rand.Intn(testSliceSize * 1 / 4)
	}

	return [][]int{
		origin[:testSliceSize/2],
		origin[testSliceSize/2:],
	}
}

func main() {
	testSlices := createTestSlices()

	// tracing
	trace.Start(os.Stderr)
	defer trace.Stop()

	// try to merge two slices into one with unique values
	uniqValues := set.NewSetInt(10)
	wg := sync.WaitGroup{}
	for i, slice := range testSlices {
		fmt.Printf("slice %d values count: %d\n", i, len(slice))
		wg.Add(1)
		go func(values []int) {
			defer wg.Done()
			for _, val := range values {
				// we use Has(val) check to reduce
				// Add(val) calls, which requires exclusive lock
				if !uniqValues.Has(val) {
					uniqValues.Add(val)
				}
			}
		}(slice)
	}
	wg.Wait()

	fmt.Printf("unique values count: %d\n", uniqValues.Len())
}
