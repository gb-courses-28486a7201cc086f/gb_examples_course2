package floatmap

import (
	"fmt"
	"sync"
	"testing"
)

var (
	setSize              = 10
	benchWorkers         = 1000
	benchWritersPctTable = []int{0, 10, 50, 90, 100}
)

type SetIface interface {
	Add(float32)
	Has(float32) bool
}

func baseSetTests(t *testing.T, set SetIface) {
	testValue := float32(5.0)

	// firstly, value not exists
	expect := false
	got := set.Has(testValue)
	if got != expect {
		t.Errorf("check not existing value failed: got %v, expect %v", got, expect)
	}
	// add value
	set.Add(testValue)
	// now it should exists
	expect = true
	got = set.Has(testValue)
	if got != expect {
		t.Errorf("check existing value failed: got %v, expect %v", got, expect)
	}
}

func TestMutexSet(t *testing.T) {
	testSet := NewMutexSet(setSize)
	baseSetTests(t, testSet)
}

func TestRwMutexSet(t *testing.T) {
	testSet := NewRWMutexSet(setSize)
	baseSetTests(t, testSet)
}

func TestSyncMapSet(t *testing.T) {
	testSet := NewSyncMapSet(setSize)
	baseSetTests(t, testSet)
}

//
// Benchmarks
//

func baseSetBench(b *testing.B, set SetIface, writePct int) {
	if writePct > 100 || writePct < 0 {
		b.Fatalf("writers %% should be int in [0, 100], got: %d", writePct)
	}
	writers := benchWorkers * writePct / 100
	readers := benchWorkers - writers

	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(benchWorkers)
		// start a few goroutines
		// readers
		for g := 0; g < readers; g++ {
			go func(i int) {
				defer wg.Done()
				set.Has(float32(i))
			}(g)
		}
		// writers
		for g := 0; g < writers; g++ {
			go func(i int) {
				defer wg.Done()
				set.Add(float32(i))
			}(g)
		}
		wg.Wait()
	}
}

func BenchmarkMutexSet(b *testing.B) {
	for _, writePct := range benchWritersPctTable {
		benchDescr := fmt.Sprintf("test set write/read with %d%% writers", writePct)
		b.Run(benchDescr, func(b *testing.B) {
			testSet := NewMutexSet(setSize)
			baseSetBench(b, testSet, writePct)
		})
	}
}

func BenchmarkRWMutexSet(b *testing.B) {
	for _, writePct := range benchWritersPctTable {
		benchDescr := fmt.Sprintf("test set write/read with %d%% writers", writePct)
		b.Run(benchDescr, func(b *testing.B) {
			testSet := NewRWMutexSet(setSize)
			baseSetBench(b, testSet, writePct)
		})
	}
}

func BenchmarkSyncMapSet(b *testing.B) {
	for _, writePct := range benchWritersPctTable {
		benchDescr := fmt.Sprintf("test set write/read with %d%% writers", writePct)
		b.Run(benchDescr, func(b *testing.B) {
			testSet := NewSyncMapSet(setSize)
			baseSetBench(b, testSet, writePct)
		})
	}
}
