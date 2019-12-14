package reader

import (
	"context"
	"errors"
	"io"
	"net/url"
	"sort"
	"strings"
	"sync"
)

var (
	readersMu sync.RWMutex
	readers   = make(map[string]Reader)
)

type Driver interface {
	Open(string) error
}

type Reader interface {
	Open(context.Context, string) error
	Read(context.Context, string) (io.ReadCloser, error)
	URI(string) string
}

func NewReader(ctx context.Context, uri string) (Reader, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	name := u.Scheme

	nrml_name := normalizeName(name)

	r, ok := readers[nrml_name]

	if !ok {
		return nil, errors.New("Unknown reader")
	}

	err = r.Open(ctx, uri)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func Register(name string, reader Reader) {

	readersMu.Lock()
	defer readersMu.Unlock()

	if reader == nil {
		panic("go-whosonfirst-reader: Register reader is nil")

	}

	nrml_name := normalizeName(name)

	if _, dup := readers[nrml_name]; dup {
		panic("go-whosonfirst-reader: Register called twice for reader " + name)
	}

	readers[nrml_name] = reader
}

func normalizeName(name string) string {
	return strings.ToUpper(name)
}

func unregisterAllReaders() {
	readersMu.Lock()
	defer readersMu.Unlock()
	readers = make(map[string]Reader)
}

func Readers() []string {

	readersMu.RLock()
	defer readersMu.RUnlock()

	var list []string

	for name := range readers {
		list = append(list, name)
	}

	sort.Strings(list)
	return list
}
