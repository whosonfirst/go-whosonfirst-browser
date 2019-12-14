package rewrite

import (
	"bufio"
	"bytes"
	"golang.org/x/net/html"
	"io"
	go_http "net/http"
	go_httptest "net/http/httptest"
	"strconv"
)

type RewriteHTMLFunc func(node *html.Node, writer io.Writer)

func RewriteHTMLHandler(prev go_http.Handler, rewrite_func RewriteHTMLFunc) go_http.Handler {

	fn := func(rsp go_http.ResponseWriter, req *go_http.Request) {

		rec := go_httptest.NewRecorder()
		prev.ServeHTTP(rec, req)

		body := rec.Body.Bytes()
		reader := bytes.NewReader(body)
		doc, err := html.Parse(reader)

		if err != nil {
			go_http.Error(rsp, err.Error(), go_http.StatusInternalServerError)
			return
		}

		var buf bytes.Buffer
		wr := bufio.NewWriter(&buf)

		rewrite_func(doc, wr)

		err = html.Render(wr, doc)

		if err != nil {
			go_http.Error(rsp, err.Error(), go_http.StatusInternalServerError)
			return
		}

		wr.Flush()

		for k, v := range rec.Header() {

			if k == "Content-Length" {
				continue
			}

			if k == "Content-Type" {
				continue
			}

			rsp.Header()[k] = v
		}

		data := buf.Bytes()
		clen := len(data)

		rsp.Header().Set("Content-Length", strconv.Itoa(clen))
		rsp.Header().Set("Content-Type", "text/html; charset=utf-8")

		rsp.WriteHeader(200)
		rsp.Write(data)
	}

	return go_http.HandlerFunc(fn)
}
