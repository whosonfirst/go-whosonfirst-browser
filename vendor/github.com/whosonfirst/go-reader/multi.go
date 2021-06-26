package reader

import (
	"context"
	"errors"
	"fmt"
	"io"
	_ "log"
	"sync"
)

type MultiReader struct {
	Reader
	readers []Reader
	lookup  map[string]int
	mu      *sync.RWMutex
}

func NewMultiReaderFromURIs(ctx context.Context, uris ...string) (Reader, error) {

	readers := make([]Reader, 0)

	for _, uri := range uris {

		r, err := NewReader(ctx, uri)

		if err != nil {
			return nil, err
		}

		readers = append(readers, r)
	}

	return NewMultiReader(ctx, readers...)
}

func NewMultiReader(ctx context.Context, readers ...Reader) (Reader, error) {

	lookup := make(map[string]int)

	mu := new(sync.RWMutex)

	mr := MultiReader{
		readers: readers,
		lookup:  lookup,
		mu:      mu,
	}

	return &mr, nil
}

func (mr *MultiReader) Read(ctx context.Context, uri string) (io.ReadSeekCloser, error) {

	missing := errors.New("Unable to read URI")

	mr.mu.RLock()

	idx, ok := mr.lookup[uri]

	mr.mu.RUnlock()

	if ok {

		// log.Printf("READ MULTIREADER LOOKUP INDEX FOR %s AS %d\n", uri, idx)

		if idx == -1 {
			return nil, missing
		}

		r := mr.readers[idx]
		return r.Read(ctx, uri)
	}

	var fh io.ReadSeekCloser
	idx = -1

	for i, r := range mr.readers {

		rsp, err := r.Read(ctx, uri)

		if err == nil {

			fh = rsp
			idx = i

			break
		}
	}

	// log.Printf("SET MULTIREADER LOOKUP INDEX FOR %s AS %d\n", uri, idx)

	mr.mu.Lock()
	mr.lookup[uri] = idx
	mr.mu.Unlock()

	if fh == nil {
		return nil, missing
	}

	return fh, nil
}

func (mr *MultiReader) ReaderURI(ctx context.Context, uri string) string {

	mr.mu.RLock()

	idx, ok := mr.lookup[uri]

	mr.mu.RUnlock()

	if ok {
		return mr.readers[idx].ReaderURI(ctx, uri)
	}

	_, err := mr.Read(ctx, uri)

	if err != nil {
		return fmt.Sprintf("x-urn:go-reader:multi#%s", uri)
	}

	return mr.ReaderURI(ctx, uri)
}
