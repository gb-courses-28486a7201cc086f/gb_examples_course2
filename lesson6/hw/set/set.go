package set

import "sync"

type SetInt struct {
	mx   sync.RWMutex
	data map[int]struct{}
}

func NewSetInt(size int) *SetInt {
	data := make(map[int]struct{}, size)
	return &SetInt{data: data}
}

func (ms *SetInt) Add(elem int) {
	ms.mx.Lock()
	defer ms.mx.Unlock()

	ms.data[elem] = struct{}{}
}

func (ms *SetInt) Has(elem int) bool {
	ms.mx.RLock()
	defer ms.mx.RUnlock()

	_, ok := ms.data[elem]
	return ok
}

func (ms *SetInt) Len() int {
	ms.mx.RLock()
	defer ms.mx.RUnlock()

	return len(ms.data)
}

func (ms *SetInt) Values() (values []int) {
	ms.mx.RLock()
	defer ms.mx.RUnlock()

	for val := range ms.data {
		values = append(values, val)
	}
	return values
}
