package reader

import (
	"context"
	"fmt"
	"io"
	_ "log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/whosonfirst/go-ioutil"
	wof_reader "github.com/whosonfirst/go-reader"
	"golang.org/x/oauth2"
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
		return nil, fmt.Errorf("Invalid path")
	}

	r.repo = parts[0]
	r.branch = DEFAULT_BRANCH

	q := u.Query()

	token := q.Get("access_token")
	branch := q.Get("branch")

	if token == "" {
		return nil, fmt.Errorf("Missing access token")
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

func (r *GitHubAPIReader) Read(ctx context.Context, uri string) (io.ReadSeekCloser, error) {

	<-r.throttle

	url := r.ReaderURI(ctx, uri)

	ref := fmt.Sprintf("refs/heads/%s", r.branch)

	opts := &github.RepositoryContentGetOptions{
		Ref: ref,
	}

	rsp, _, _, err := r.client.Repositories.GetContents(ctx, r.owner, r.repo, url, opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to get contents for %s, %w", url, err)
	}

	body, err := rsp.GetContent()

	var rsp_r io.Reader

	if err != nil {

		// START OF I have no idea why I need to do this, but only sometimes...

		if *rsp.Content != "" {
			return nil, fmt.Errorf("Failed to read contents for %s, %w", url, err)
		}

		if *rsp.DownloadURL == "" {
			return nil, fmt.Errorf("Failed to read contents for %s and response is missing download URL, %w", url, err)
		}

		raw_rsp, err := http.Get(*rsp.DownloadURL)

		if err != nil {
			return nil, fmt.Errorf("Failed to fetch contents from download URL, %w", err)
		}

		rsp_r = raw_rsp.Body

		// END OF I have no idea why I need to do this, but only sometimes...

	} else {
		rsp_r = strings.NewReader(body)
	}

	fh, err := ioutil.NewReadSeekCloser(rsp_r)

	if err != nil {
		return nil, fmt.Errorf("Failed to create ReadSeekCloser for %s, %w", url, err)
	}

	return fh, nil
}

func (r *GitHubAPIReader) ReaderURI(ctx context.Context, key string) string {

	uri := key

	if r.prefix != "" {
		uri = filepath.Join(r.prefix, key)
	}

	return uri
}
