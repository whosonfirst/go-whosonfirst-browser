package service

import (
	"context"
	"errors"
	"github.com/aaronland/go-artisanal-integers"
	"github.com/aaronland/go-pool"
	"github.com/whosonfirst/go-whosonfirst-log"
	"math/rand"
	"sync"
	"time"
)

var r *rand.Rand

func init() {

	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type ProxyServiceOptions struct {
	Pool    pool.Pool
	Minimum int
	Logger  *log.WOFLogger
	Workers int
}

func DefaultProxyServiceOptions() (*ProxyServiceOptions, error) {

	ctx := context.Background()
	pl, err := pool.NewPool(ctx, "memory://")

	if err != nil {
		return nil, err
	}

	logger := log.NewWOFLogger()

	opts := ProxyServiceOptions{
		Pool:    pl,
		Logger:  logger,
		Minimum: 10,
		Workers: 0,
	}

	return &opts, nil
}

type ProxyService struct {
	artisanalinteger.Service
	artisanalinteger.Client
	options *ProxyServiceOptions
	clients []artisanalinteger.Client
	refill  chan bool
}

func NewProxyService(opts *ProxyServiceOptions, clients ...artisanalinteger.Client) (artisanalinteger.Service, error) {

	if len(clients) == 0 {
		return nil, errors.New("Insuffient clients")
	}

	// See notes in RefillPool() for details

	size := opts.Minimum
	refill := make(chan bool, size)

	for i := 0; i < size; i++ {
		refill <- true
	}

	// Possibly also keep global stats on number of fetches
	// and cache hits/misses/etc

	p := &ProxyService{
		options: opts,
		clients: clients,
		refill:  refill,
	}

	go p.refillPool()
	go p.status()
	go p.monitor()

	return p, nil
}

func (p *ProxyService) status() {

	for {
		select {
		case <-time.After(5 * time.Second):
			p.options.Logger.Status("pool length: %d", p.options.Pool.Length())
		}
	}
}

func (p *ProxyService) monitor() {

	for {
		select {
		case <-time.After(10 * time.Second):
			if p.options.Pool.Length() < int64(p.options.Minimum) {
				go p.refillPool()
			}
		}

	}
}

func (p *ProxyService) refillPool() {

	// Remember there is a fixed size work queue of allowable times to try
	// and refill the pool simultaneously. First, we block until a slot opens
	// up.

	<-p.refill

	t1 := time.Now()

	// Figure out how many integers we need to get *at this moment* which when
	// the service is under heavy load is a misleading number at best. It might
	// be worth adjusting this by a factor of (n) depending on the current load.
	// But that also means tracking what we think the current load means so we
	// aren't going to do that now...

	todo := int64(p.options.Minimum) - p.options.Pool.Length()

	workers := p.options.Workers

	if workers == 0 {
		workers = int(p.options.Minimum / 2)
	}

	if workers == 0 {
		workers = 1
	}

	// Now we're going to set up two simultaneous queues. One (the work group) is
	// just there to keep track of all the requests for new integers we need to
	// make. The second (the throttle) is there to make sure we don't exhaust all
	// the filehandles or network connections.

	th := make(chan bool, workers)

	for i := 0; i < workers; i++ {
		th <- true
	}

	wg := new(sync.WaitGroup)

	p.options.Logger.Debug("refill poll w/ %d integers and %d workers", todo, workers)

	success := 0
	failed := 0

	for j := 0; int64(j) < todo; j++ {

		// Wait for the throttle to open a slot. Also record whether
		// the operation was successful.

		rsp := <-th

		if rsp == true {
			success += 1
		} else {
			failed += 1
		}

		// First check that we still actually need to keep fetching integers

		if p.options.Pool.Length() >= int64(p.options.Minimum) {
			p.options.Logger.Debug("pool is full (%d) stopping after %d iterations", p.options.Pool.Length(), j)
			break
		}

		// Standard work group stuff

		wg.Add(1)

		// Sudo make me a sandwitch. Note the part where we ping the throttle with
		// the return value at the end both to signal an available slot and to record
		// whether the integer harvesting was successful.

		go func(pr *ProxyService) {
			defer wg.Done()
			th <- pr.addToPool()
		}(p)
	}

	// More standard work group stuff

	wg.Wait()

	// Again note the way we are freeing a spot in the refill queue

	p.refill <- true

	t2 := time.Since(t1)
	p.options.Logger.Info("time to refill the pool with %d integers (success: %d failed: %d): %v (pool length is now %d)", todo, success, failed, t2, p.options.Pool.Length())

}

func (p *ProxyService) addToPool() bool {

	i, err := p.getInteger()

	if err != nil {
		return false
	}

	pi := pool.NewIntItem(i)

	p.options.Pool.Push(pi)
	return true
}

func (p *ProxyService) getInteger() (int64, error) {

	idx := r.Intn(len(p.clients))

	cl := p.clients[idx] // round-robin me or something...

	i, err := cl.NextInt()

	if err != nil {
		p.options.Logger.Error("failed to create new integer, because %v", err)
		return 0, err
	}

	p.options.Logger.Debug("got new integer %d", i)
	return i, nil
}

func (p *ProxyService) NextInt() (int64, error) {

	if p.options.Pool.Length() == 0 {

		go p.refillPool()

		p.options.Logger.Warning("pool length is 0 so fetching integer from source")
		return p.getInteger()
	}

	v, ok := p.options.Pool.Pop()

	if !ok {
		p.options.Logger.Error("failed to pop integer!")

		go p.refillPool()
		return p.getInteger()
	}

	i := v.Int()

	p.options.Logger.Debug("return cached integer %d", i)

	return i, nil
}

func (p *ProxyService) LastInt() (int64, error) {
	return -1, errors.New("Please write me")
}
