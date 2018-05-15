package buffalloc

import (
	"fmt"
	"sync"
)

var (
	invalidSizeError = fmt.Errorf("invalid buffer size")
)

// BuffAllocator is used to allocate and release buffers efficiently without large GC overhead.
//
type BuffAllocator interface {
	// Allocate allocates a new buffer whose capacity is larger than or equal to the specified size
	// Allocate assures that cap(buff) >= size and len(buff) == 0
	Allocate(size int) (buff []byte)
	// Release releases a buffer which must be previously allocated by Allocate of the same BuffAllocator
	Release(buff []byte)
}

// NewBuffAllocator creates a new buffer allocator
func NewBuffAllocator(minimalSize int, multiply int) BuffAllocator {
	if minimalSize < 1 {
		panic(fmt.Errorf("minimalSize should be at least 1"))
	}

	if multiply < 2 {
		panic(fmt.Errorf("multiply should be at least 2"))
	}

	var ba BuffAllocator
	ba = &buffAllocator{
		minimalSize: minimalSize,
		multiply:    multiply,
		pools:       nil,
	}
	return ba
}

type buffAllocator struct {
	minimalSize int
	multiply    int
	pools       []*sync.Pool
}

func (ba *buffAllocator) Allocate(size int) (buff []byte) {
	pool := ba.getPool(size)
	return pool.Get().([]byte)
}

func (ba *buffAllocator) Release(buff []byte) {
	buffcap := cap(buff)
	ba.validateReleaseBuffSize(buffcap)
	buff = buff[:0]
	ba.getPool(buffcap).Put(buff)
}

func (ba *buffAllocator) validateReleaseBuffSize(size int) {
	normSize := ba.minimalSize
	for normSize < size {
		normSize *= ba.multiply
	}
	if normSize != size {
		panic(invalidSizeError)
	}
}

func (ba *buffAllocator) getPool(size int) (pool *sync.Pool) {
	normSize := ba.minimalSize
	poolIdx := 0
	for normSize < size {
		normSize *= ba.multiply
		poolIdx += 1
	}
	for len(ba.pools) < poolIdx+1 {
		ba.pools = append(ba.pools, nil)
	}
	pool = ba.pools[poolIdx]
	if pool == nil {
		pool = &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, normSize)
			},
		}
		ba.pools[poolIdx] = pool
	}
	return
}
