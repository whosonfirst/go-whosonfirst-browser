package reader

import (
	"context"
	"errors"
	wof_reader "github.com/whosonfirst/go-reader"
	"io"
	_ "log"
	"net/http"
	"net/url"
	"path/filepath"
	"time"
)

func init() {

	ctx := context.Background()

	schemes := []string{
		"http",
		"https",			
	}

	for _, s := range schemes {
	
		err := wof_reader.RegisterReader(ctx, s, initializeHTTPReader)	

		if err != nil {
			panic(err)
		}
	}
}

func initializeHTTPReader(ctx context.Context, uri string) (wof_reader.Reader, error) {

	r := NewHTTPReader()
	err := r.Open(ctx, uri)

	if err != nil {
		return nil, err
	}

	return r, nil
}

type HTTPReader struct {
	wof_reader.Reader
	url *url.URL
	throttle <-chan time.Time
}

func NewHTTPReader() wof_reader.Reader {

	rate := time.Second / 3
	throttle := time.Tick(rate)

	r := HTTPReader{
		throttle: throttle,
	}

	return &r
}

func (r *HTTPReader) Open(ctx context.Context, uri string) error {

	u, err := url.Parse(uri)

	if err != nil {
		return err
	}

	r.url = u
	return nil
}

func (r *HTTPReader) Read(ctx context.Context, uri string) (io.ReadCloser, error) {

	<-r.throttle

	u, _ := url.Parse(r.url.String())
	u.Path = filepath.Join(u.Path, uri)
	
	url := u.String()

	rsp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != 200 {
		return nil, errors.New(rsp.Status)
	}

	return rsp.Body, nil
}
