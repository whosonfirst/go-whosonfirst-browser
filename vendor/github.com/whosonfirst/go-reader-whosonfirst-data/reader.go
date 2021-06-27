package reader

import (
	"context"
	"encoding/json"
	"fmt"
	wof_reader "github.com/whosonfirst/go-reader"
	_ "github.com/whosonfirst/go-reader-github"
	wof_uri "github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type WhosOnFirstDataReader struct {
	wof_reader.Reader
	throttle     <-chan time.Time
	provider     string
	organization string
	repo         string
	repos        *sync.Map
	readers      *sync.Map
}

func init() {

	ctx := context.Background()
	err := wof_reader.RegisterReader(ctx, "whosonfirst-data", NewWhosOnFirstDataReader)

	if err != nil {
		panic(err)
	}
}

func NewWhosOnFirstDataReader(ctx context.Context, uri string) (wof_reader.Reader, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	provider := q.Get("provider")
	org := q.Get("organization")
	repo := q.Get("repo")

	if provider == "" {
		provider = "github"
	}

	if org == "" {
		org = "whosonfirst-data"
	}

	rate := time.Second / 3
	throttle := time.Tick(rate)

	repos := new(sync.Map)
	readers := new(sync.Map)

	r := &WhosOnFirstDataReader{
		throttle:     throttle,
		provider:     provider,
		organization: org,
		repo:         repo,
		repos:        repos,
		readers:      readers,
	}

	return r, nil
}

func (r *WhosOnFirstDataReader) Read(ctx context.Context, uri string) (io.ReadSeekCloser, error) {

	id, _, err := wof_uri.ParseURI(uri)

	if err != nil {
		return nil, err
	}

	<-r.throttle

	select {
	case <-ctx.Done():
		return nil, nil
	default:
		// pass
	}

	repo := r.repo

	if repo == "" {

		this_repo, err := r.getRepo(ctx, id)

		if err != nil {
			return nil, err
		}

		repo = this_repo
	}

	gh_r, err := r.getReader(ctx, repo)

	if err != nil {
		return nil, err
	}

	return gh_r.Read(ctx, uri)
}

func (r *WhosOnFirstDataReader) getReader(ctx context.Context, repo string) (wof_reader.Reader, error) {

	v, ok := r.readers.Load(repo)

	if ok {
		gh_r := v.(wof_reader.Reader)
		return gh_r, nil
	}

	gh_uri := fmt.Sprintf("%s://%s/%s", r.provider, r.organization, repo)

	gh_r, err := wof_reader.NewReader(ctx, gh_uri)

	if err != nil {
		return nil, err
	}

	go func() {
		r.readers.Store(repo, gh_r)
	}()

	return gh_r, nil
}

func (r *WhosOnFirstDataReader) getRepo(ctx context.Context, id int64) (string, error) {

	v, ok := r.repos.Load(id)

	if ok {
		repo := v.(string)
		return repo, nil
	}

	uri := fmt.Sprintf("https://data.whosonfirst.org/findingaid/%d", id)

	rsp, err := http.Get(uri)

	if err != nil {
		return "", err
	}

	defer rsp.Body.Close()

	// https://github.com/whosonfirst/go-whosonfirst-findingaid/blob/master/repo/repo.go

	type FindingAidResponse struct {
		ID   int64  `json:"id"`
		Repo string `json:"repo"`
		URI  string `json:"uri"`
	}

	var fa_rsp *FindingAidResponse

	dec := json.NewDecoder(rsp.Body)

	err = dec.Decode(&fa_rsp)

	if err != nil {
		return "", err
	}

	repo := fa_rsp.Repo

	go func() {
		r.repos.Store(id, repo)
	}()

	return repo, nil
}
