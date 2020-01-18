package pool

// https://github.com/SimonWaldherr/golang-examples/blob/2be89f3185aded00740a45a64e3c98855193b948/advanced/lifo.go

import (
	"context"
	"sync"
	"sync/atomic"
)

func init() {
	ctx := context.Background()
	pl := NewMemoryPool()
	Register(ctx, "memory", pl)
}

type MemoryPool struct {
	Pool
	nodes []Item
	count int64
	mutex *sync.Mutex
}

func NewMemoryPool() Pool {

	mu := new(sync.Mutex)
	nodes := make([]Item, 0)

	pl := &MemoryPool{
		mutex: mu,
		nodes: nodes,
		count: 0,
	}

	return pl
}

func (pl *MemoryPool) Open(ctx context.Context, uri string) error {
	return nil
}

func (pl *MemoryPool) Length() int64 {

	pl.mutex.Lock()
	defer pl.mutex.Unlock()

	return atomic.LoadInt64(&pl.count)
}

func (pl *MemoryPool) Push(i Item) {

	pl.mutex.Lock()
	defer pl.mutex.Unlock()

	pl.nodes = append(pl.nodes[:pl.count], i)
	atomic.AddInt64(&pl.count, 1)
}

func (pl *MemoryPool) Pop() (Item, bool) {

	pl.mutex.Lock()
	defer pl.mutex.Unlock()

	if pl.count == 0 {
		return nil, false
	}

	atomic.AddInt64(&pl.count, -1)
	i := pl.nodes[pl.count]

	return i, true
}
