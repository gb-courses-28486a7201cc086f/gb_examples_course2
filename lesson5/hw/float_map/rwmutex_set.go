// Package floatmap contains a few
// Set datastucture implementations
// using different sync primitives
package floatmap

import "sync"

type RWMutexSet struct {
	mx   sync.RWMutex
	data map[float32]struct{}
}

func NewRWMutexSet(size int) *RWMutexSet {
	data := make(map[float32]struct{}, size)
	return &RWMutexSet{data: data}
}

func (ms *RWMutexSet) Add(elem float32) {
	ms.mx.Lock()
	defer ms.mx.Unlock()

	ms.data[elem] = struct{}{}
}

func (ms *RWMutexSet) Has(elem float32) bool {
	ms.mx.RLock()
	defer ms.mx.RUnlock()

	_, ok := ms.data[elem]
	return ok
}
