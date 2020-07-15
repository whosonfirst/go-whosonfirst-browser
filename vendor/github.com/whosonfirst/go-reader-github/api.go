package reader

import (
	"context"
	"errors"
	"github.com/google/go-github/github"
	wof_reader "github.com/whosonfirst/go-reader"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	_ "log"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

type GitHubAPIReader struct {
	wof_reader.Reader
	owner    string
	repo     string
	prefix   string
	branch   string
	client   *github.Client
	throttle <-chan time.Time
}

func init() {

	ctx := context.Background()
	err := wof_reader.RegisterReader(ctx, "githubapi", NewGitHubAPIReader)

	if err != nil {
		panic(err)
	}
}

func NewGitHubAPIReader(ctx context.Context, uri string) (wof_reader.Reader, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	rate := time.Second / 3
	throttle := time.Tick(rate)

	r := &GitHubAPIReader{
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

	token := q.Get("access_token")
	branch := q.Get("branch")

	if token == "" {
		return nil, errors.New("Missing access token")
	}

	if branch != "" {
		r.branch = branch
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	r.client = client

	prefix := q.Get("prefix")
	r.prefix = prefix

	return r, nil
}

func (r *GitHubAPIReader) Read(ctx context.Context, uri string) (io.ReadCloser, error) {

	<-r.throttle

	url := r.URI(uri)

	opts := &github.RepositoryContentGetOptions{}

	rsp, _, _, err := r.client.Repositories.GetContents(ctx, r.owner, r.repo, url, opts)

	if err != nil {
		return nil, err
	}

	body, err := rsp.GetContent()

	if err != nil {
		return nil, err
	}

	br := strings.NewReader(body)
	fh := ioutil.NopCloser(br)

	return fh, nil
}

func (r *GitHubAPIReader) URI(key string) string {

	uri := key

	if r.prefix != "" {
		uri = filepath.Join(r.prefix, key)
	}

	return uri
}
