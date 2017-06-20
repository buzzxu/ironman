package ironman

import "sync"

// DEFAULT_SYNC_POOL 默認對象池
var DEFAULT_SYNC_POOL *SyncPool

// NewPool is 分配内存的时候，从池子里面找满足容量切最小的池子。比如申请长度是2的，就分配大小为5的那个池子。如果是11，就分配大小是20的那个池子里面的对象
func NewPool() *SyncPool {
	DEFAULT_SYNC_POOL = NewSyncPool(
		5,
		30000,
		2,
	)
	return DEFAULT_SYNC_POOL
}

// Alloc is 分配大小
func Alloc(size int) []int64 {
	return DEFAULT_SYNC_POOL.Alloc(size)
}

// Free 釋放内存
func Free(mem []int64) {
	DEFAULT_SYNC_POOL.Free(mem)
}

// SyncPool is a sync.Pool base slab allocation memory pool
type SyncPool struct {
	classes     []sync.Pool
	classesSize []int
	minSize     int
	maxSize     int
}

// NewSyncPool 申請
func NewSyncPool(minSize, maxSize, factor int) *SyncPool {
	n := 0
	for chunkSize := minSize; chunkSize <= maxSize; chunkSize *= factor {
		n++
	}
	pool := &SyncPool{
		make([]sync.Pool, n),
		make([]int, n),
		minSize, maxSize,
	}
	n = 0
	for chunkSize := minSize; chunkSize <= maxSize; chunkSize *= factor {
		pool.classesSize[n] = chunkSize
		pool.classes[n].New = func(size int) func() interface{} {
			return func() interface{} {
				buf := make([]int64, size)
				return &buf
			}
		}(chunkSize)
		n++
	}
	return pool
}

// Alloc 申請内存
func (pool *SyncPool) Alloc(size int) []int64 {
	if size <= pool.maxSize {
		for i := 0; i < len(pool.classesSize); i++ {
			if pool.classesSize[i] >= size {
				mem := pool.classes[i].Get().(*[]int64)
				// return (*mem)[:size]
				return (*mem)[:0]
			}
		}
	}
	return make([]int64, 0, size)
}

// Free 釋放
func (pool *SyncPool) Free(mem []int64) {
	if size := cap(mem); size <= pool.maxSize {
		for i := 0; i < len(pool.classesSize); i++ {
			if pool.classesSize[i] >= size {
				pool.classes[i].Put(&mem)
				return
			}
		}
	}
}
