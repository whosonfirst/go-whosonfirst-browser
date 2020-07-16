package reader

import (
	"context"
	"errors"
	"fmt"
	wof_reader "github.com/whosonfirst/go-reader"
	"io"
	_ "log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type GitHubReader struct {
	wof_reader.Reader
	owner    string
	repo     string
	branch   string
	throttle <-chan time.Time
}

func init() {

	ctx := context.Background()
	err := wof_reader.RegisterReader(ctx, "github", NewGitHubReader)

	if err != nil {
		panic(err)
	}
}

func NewGitHubReader(ctx context.Context, uri string) (wof_reader.Reader, error) {

	rate := time.Second / 3
	throttle := time.Tick(rate)

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	r := &GitHubReader{
		throttle: throttle,
	}

	r.owner = u.Host

	path := strings.TrimLeft(u.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) != 1 {
		return nil, errors.New("Invalid path")
	}

	r.repo = parts[0]
	r.branch = "master"

	q := u.Query()

	branch := q.Get("branch")

	if branch != "" {
		r.branch = branch
	}

	return r, nil
}

func (r *GitHubReader) Read(ctx context.Context, uri string) (io.ReadCloser, error) {

	<-r.throttle

	url := r.URI(uri)

	rsp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != 200 {
		return nil, errors.New(rsp.Status)
	}

	return rsp.Body, nil
}

func (r *GitHubReader) URI(key string) string {

	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/data/%s", r.owner, r.repo, r.branch, key)
}
