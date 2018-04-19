package reader

import (
	wof_reader "github.com/whosonfirst/go-whosonfirst-readwrite/reader"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type HTTPReader struct {
	wof_reader.Reader
	root *url.URL
}

func NewHTTPReader(root string) (wof_reader.Reader, error) {

	root_url, err := url.Parse(root)

	if err != nil {
		return nil, err
	}

	r := HTTPReader{
		root: root_url,
	}

	return &r, nil
}

func (r *HTTPReader) Read(path string) (io.ReadCloser, error) {

     	uri := r.URI(path)

	rsp, err := http.Get(uri)

	if err != nil {
		return nil, err
	}

	return rsp.Body, nil
}

func (r *HTTPReader) URI(path string) string {

	url := r.root.String() + path

	if !strings.HasSuffix(r.root.String(), "/") && !strings.HasPrefix(path, "/") {
		url = r.root.String() + "/" + path
	}

	return url
}
