package reader

import (
	"io"
	_ "log"		
	"net/http"
	"net/url"
)

type HTTPReader struct {
	Reader
	root *url.URL
}

func NewHTTPReader(root *url.URL) (Reader, error) {

	r := HTTPReader{
		root: root,
	}

	return &r, nil
}

func (r *HTTPReader) Read(key string) (io.ReadCloser, error) {

	url := r.root.String() + key
	// log.Println("FETCH", url)
	
	rsp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	return rsp.Body, nil
}
