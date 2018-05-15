package go_buff_alloc

import (
	"fmt"
	"sync"
)

var (
	invalidSizeError = fmt.Errorf("invalid buffer size")
)

type BuffAllocator interface {
	Allocate(size int) (buff []byte)
	Release(buff []byte)
}

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
	buff = buff[:buffcap]
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
				return make([]byte, normSize)
			},
		}
		ba.pools[poolIdx] = pool
	}
	return
}
