package writer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"sort"
	"strings"
	"sync"
)

var (
	writersMu sync.RWMutex
	writers   = make(map[string]Writer)
)

type Driver interface {
	Open(string) error
}

type Writer interface {
	Open(context.Context, string) error
	Write(context.Context, string, io.ReadCloser) error
	URI(string) string
}

func NewWriter(ctx context.Context, uri string) (Writer, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	name := u.Scheme

	nrml_name := normalizeName(name)

	r, ok := writers[nrml_name]

	if !ok {
		return nil, errors.New(fmt.Sprintf("Unknown writer '%s' (normalized as '%s')", name, nrml_name))
	}

	err = r.Open(ctx, uri)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func Register(name string, writer Writer) {

	writersMu.Lock()
	defer writersMu.Unlock()

	if writer == nil {
		panic("go-whosonfirst-Writer: Register writer is nil")

	}

	nrml_name := normalizeName(name)

	if _, dup := writers[nrml_name]; dup {
		panic("go-whosonfirst-writer: Register called twice for writer " + name)
	}

	writers[nrml_name] = writer
}

func normalizeName(name string) string {
	return strings.ToUpper(name)
}

func unregisterAllWriters() {
	writersMu.Lock()
	defer writersMu.Unlock()
	writers = make(map[string]Writer)
}

func Writers() []string {

	writersMu.RLock()
	defer writersMu.RUnlock()

	var list []string

	for name := range writers {
		list = append(list, name)
	}

	sort.Strings(list)
	return list
}
