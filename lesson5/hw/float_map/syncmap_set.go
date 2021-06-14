// Package floatmap contains a few
// Set datastucture implementations
// using different sync primitives
package floatmap

import "sync"

type SyncMapSet struct {
	data sync.Map
}

func NewSyncMapSet(size int) *SyncMapSet {
	return &SyncMapSet{}
}

func (ms *SyncMapSet) Add(elem float32) {
	ms.data.Store(elem, struct{}{})
}

func (ms *SyncMapSet) Has(elem float32) bool {
	_, ok := ms.data.Load(elem)
	return ok
}
