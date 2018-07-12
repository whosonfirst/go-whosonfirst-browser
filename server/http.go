package server

import (
	_ "log"
	"net/http"
	"net/url"
)

type HTTPServer struct {
	Server
	url *url.URL
}

func NewHTTPServer(u *url.URL, args ...interface{}) (Server, error) {

	u.Scheme = "http"

	server := HTTPServer{
		url: u,
	}

	return &server, nil
}

func (s *HTTPServer) Address() string {
	return s.url.String()
}

func (s *HTTPServer) ListenAndServe(mux *http.ServeMux) error {
	return http.ListenAndServe(s.url.Host, mux)
}
