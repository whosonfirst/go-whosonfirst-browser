package iterator

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-iterate/emitter"
	"io"
	"log"
	"net/url"
	"regexp"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

type Iterator struct {
	Emitter             emitter.Emitter
	EmitterCallbackFunc emitter.EmitterCallbackFunc
	Logger              *log.Logger
	Seen                int64
	count               int64
	max_procs           int
	exclude_paths       *regexp.Regexp
}

func NewIterator(ctx context.Context, emitter_uri string, emitter_cb emitter.EmitterCallbackFunc) (*Iterator, error) {

	idx, err := emitter.NewEmitter(ctx, emitter_uri)

	if err != nil {
		return nil, err
	}

	u, err := url.Parse(emitter_uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	max_procs := runtime.NumCPU()

	if q.Get("_max_procs") != "" {

		max, err := strconv.ParseInt(q.Get("_max_procs"), 10, 64)

		if err != nil {
			return nil, err
		}

		max_procs = int(max)
	}

	logger := log.Default()

	i := Iterator{
		Emitter:             idx,
		EmitterCallbackFunc: emitter_cb,
		Logger:              logger,
		Seen:                0,
		count:               0,
		max_procs:           max_procs,
	}

	if q.Get("_exclude") != "" {

		re_exclude, err := regexp.Compile(q.Get("_exclude"))

		if err != nil {
			return nil, err
		}

		i.exclude_paths = re_exclude
	}

	return &i, nil
}

func (idx *Iterator) IterateURIs(ctx context.Context, uris ...string) error {

	t1 := time.Now()

	defer func() {
		t2 := time.Since(t1)
		idx.Logger.Printf("time to index paths (%d) %v", len(uris), t2)
	}()

	idx.increment()
	defer idx.decrement()

	local_callback := func(ctx context.Context, fh io.ReadSeeker, args ...interface{}) error {

		defer atomic.AddInt64(&idx.Seen, 1)

		if idx.exclude_paths != nil {

			path, err := emitter.PathForContext(ctx)

			if err != nil {
				return err
			}

			if idx.exclude_paths.MatchString(path) {
				return nil
			}
		}

		return idx.EmitterCallbackFunc(ctx, fh, args...)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	procs := idx.max_procs
	throttle := make(chan bool, procs)

	for i := 0; i < procs; i++ {
		throttle <- true
	}

	done_ch := make(chan bool)
	err_ch := make(chan error)

	remaining := len(uris)

	for _, uri := range uris {

		go func(uri string) {

			<-throttle

			defer func() {
				throttle <- true
				done_ch <- true
			}()

			select {
			case <-ctx.Done():
				return
			default:
				// pass
			}

			err := idx.Emitter.WalkURI(ctx, local_callback, uri)

			if err != nil {
				err_ch <- err
			}
		}(uri)
	}

	for remaining > 0 {
		select {
		case <-done_ch:
			remaining -= 1
		case err := <-err_ch:
			return err
		default:
			// pass
		}
	}

	return nil
}

func (idx *Iterator) IsIndexing() bool {

	if atomic.LoadInt64(&idx.count) > 0 {
		return true
	}

	return false
}

func (idx *Iterator) increment() {
	atomic.AddInt64(&idx.count, 1)
}

func (idx *Iterator) decrement() {
	atomic.AddInt64(&idx.count, -1)
}
