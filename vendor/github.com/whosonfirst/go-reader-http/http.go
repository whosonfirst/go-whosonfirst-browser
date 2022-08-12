package reader

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-ioutil"
	wof_reader "github.com/whosonfirst/go-reader"
	"io"
	_ "log"
	"net/http"
	"net/url"
	"path/filepath"
	"time"
)

type HTTPReader struct {
	wof_reader.Reader
	url        *url.URL
	throttle   <-chan time.Time
	user_agent string
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

	q := u.Query()
	ua := q.Get("user-agent")

	if ua != "" {
		r.user_agent = ua
	}

	return &r, nil
}

func (r *HTTPReader) Read(ctx context.Context, uri string) (io.ReadSeekCloser, error) {

	<-r.throttle

	u, _ := url.Parse(r.url.String())
	u.Path = filepath.Join(u.Path, uri)

	url := u.String()

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new request, %w", err)
	}

	if r.user_agent != "" {
		req.Header.Set("User-Agent", r.user_agent)
	}

	cl := &http.Client{}

	rsp, err := cl.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Failed to execute request, %w", err)
	}

	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("Unexpected status code: %s", rsp.Status)
	}

	fh, err := ioutil.NewReadSeekCloser(rsp.Body)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new ReadSeekCloser, %w", err)
	}

	return fh, nil
}

func (r *HTTPReader) ReaderURI(ctx context.Context, uri string) string {
	return uri
}
