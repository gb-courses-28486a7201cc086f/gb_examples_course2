// Package floatmap contains a few
// Set datastucture implementations
// using different sync primitives
package floatmap

import "sync"

type MutexSet struct {
	mx   sync.Mutex
	data map[float32]struct{}
}

func NewMutexSet(size int) *MutexSet {
	data := make(map[float32]struct{}, size)
	return &MutexSet{data: data}
}

func (ms *MutexSet) Add(elem float32) {
	ms.mx.Lock()
	defer ms.mx.Unlock()

	ms.data[elem] = struct{}{}
}

func (ms *MutexSet) Has(elem float32) bool {
	ms.mx.Lock()
	defer ms.mx.Unlock()

	_, ok := ms.data[elem]
	return ok
}
