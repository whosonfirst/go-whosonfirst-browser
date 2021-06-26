package reader

import (
	"context"
	"errors"
	wof_reader "github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-ioutil"
	"io"
	_ "log"
	"net/http"
	"net/url"
	"path/filepath"
	"time"
)

type HTTPReader struct {
	wof_reader.Reader
	url      *url.URL
	throttle <-chan time.Time
}

func init() {

	ctx := context.Background()

	schemes := []string{
		"http",
		"https",
	}

	for _, s := range schemes {

		err := wof_reader.RegisterReader(ctx, s, NewHTTPReader)

		if err != nil {
			panic(err)
		}
	}
}

func NewHTTPReader(ctx context.Context, uri string) (wof_reader.Reader, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	rate := time.Second / 3
	throttle := time.Tick(rate)

	r := HTTPReader{
		throttle: throttle,
		url:      u,
	}

	return &r, nil
}

func (r *HTTPReader) Read(ctx context.Context, uri string) (io.ReadSeekCloser, error) {

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

	fh, err := ioutil.NewReadSeekCloser(rsp.Body)

	if err != nil {
		return nil, err
	}

	return fh, nil
}

func (r *HTTPReader) ReaderURI(ctx context.Context, uri string) string {
	return uri
}
