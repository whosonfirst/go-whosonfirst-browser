package writer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/go-github/v48/github"
	wof_writer "github.com/whosonfirst/go-writer/v3"
	"golang.org/x/oauth2"
)

const GITHUBAPI_SCHEME string = "githubapi"

type GitHubAPIWriterCommitTemplates struct {
	New    string
	Update string
}

type GitHubAPIWriter struct {
	wof_writer.Writer
	owner              string
	repo               string
	branch             string
	prefix             string
	client             *github.Client
	user               *github.User
	throttle           <-chan time.Time
	templates          *GitHubAPIWriterCommitTemplates
	retry_on_ratelimit bool
	retry_on_conflict  bool
	retry_attempts     int32
	max_retry_attempts int32
}

func init() {

	ctx := context.Background()
	wof_writer.RegisterWriter(ctx, GITHUBAPI_SCHEME, NewGitHubAPIWriter)
}

func NewGitHubAPIWriter(ctx context.Context, uri string) (wof_writer.Writer, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	path := strings.TrimLeft(u.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) != 1 {
		return nil, errors.New("Invalid path")
	}

	repo := parts[0]
	branch := DEFAULT_BRANCH

	q := u.Query()

	token := q.Get("access_token")

	prefix := q.Get("prefix")
	q_branch := q.Get("branch")

	if token == "" {
		return nil, errors.New("Missing access token")
	}

	if q_branch != "" {
		branch = q_branch
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	users := client.Users
	user, _, err := users.Get(ctx, "")

	if err != nil {
		return nil, err
	}

	new_template := q.Get("new")
	update_template := q.Get("update")

	if new_template == "" {
		new_template = "Created %s"
	}

	if update_template == "" {
		update_template = "Updated %s"
	}

	templates := &GitHubAPIWriterCommitTemplates{
		New:    new_template,
		Update: update_template,
	}

	retry_on_ratelimit := false
	str_ratelimit := q.Get("retry-on-ratelimit")

	if str_ratelimit != "" {

		r, err := strconv.ParseBool(str_ratelimit)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse retry-on-ratelimit parameter, %w", err)
		}

		retry_on_ratelimit = r
	}

	retry_on_conflict := false
	str_conflict := q.Get("retry-on-conflict")

	if str_conflict != "" {

		r, err := strconv.ParseBool(str_conflict)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse retry-on-conflict parameter, %w", err)
		}

		retry_on_conflict = r
	}

	max_retries := int32(10)
	str_retries := q.Get("max-retry-attempts")

	if str_retries != "" {

		r, err := strconv.Atoi(str_retries)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse max-retry-attempts parameter, %w", err)
		}

		if r < 0 {
			r = 0
		}

		max_retries = int32(r)
	}

	rate := time.Second / 3
	throttle := time.Tick(rate)

	wr := &GitHubAPIWriter{
		client:             client,
		owner:              u.Host,
		user:               user,
		repo:               repo,
		branch:             branch,
		prefix:             prefix,
		templates:          templates,
		throttle:           throttle,
		retry_on_ratelimit: retry_on_ratelimit,
		retry_on_conflict:  retry_on_conflict,
		max_retry_attempts: max_retries,
	}

	return wr, nil
}

func (wr *GitHubAPIWriter) Write(ctx context.Context, uri string, fh io.ReadSeeker) (int64, error) {

	<-wr.throttle

	body, err := io.ReadAll(fh)

	if err != nil {
		return 0, err
	}

	url := wr.WriterURI(ctx, uri)

	commit_msg := fmt.Sprintf(wr.templates.New, url)
	name := *wr.user.Login
	email := fmt.Sprintf("%s@localhost", name)

	update_opts := &github.RepositoryContentFileOptions{
		Message: github.String(commit_msg),
		Content: body,
		Branch:  github.String(wr.branch),
		Committer: &github.CommitAuthor{
			Name:  github.String(name),
			Email: github.String(email),
		},
	}

	get_opts := &github.RepositoryContentGetOptions{
		Ref: wr.branch,
	}

	get_rsp, _, _, err := wr.client.Repositories.GetContents(ctx, wr.owner, wr.repo, url, get_opts)

	if err == nil {
		commit_msg = fmt.Sprintf(wr.templates.Update, url)
		update_opts.Message = github.String(commit_msg)
		update_opts.SHA = get_rsp.SHA
	}

	_, update_rsp, err := wr.client.Repositories.UpdateFile(ctx, wr.owner, wr.repo, url, update_opts)

	if err != nil {

		try_to_recover := false

		if update_rsp.StatusCode == 409 && wr.retry_on_conflict {
			try_to_recover = true
		}

		ratelimit_err, is_ratelimit := err.(*github.RateLimitError)

		if is_ratelimit && wr.retry_on_ratelimit {
			try_to_recover = true
		}

		if !try_to_recover {
			return 0, fmt.Errorf("Failed to update %s, %w", url, err)
		}

		_, err = fh.Seek(0, 0)

		if err != nil {
			return 0, fmt.Errorf("Trigger a rate limit error but unable to rewind filehandle, %w", err)
		}

		// Try not to spin madly out of control

		if wr.max_retry_attempts > 0 {

			atomic.AddInt32(&wr.retry_attempts, 1)

			if atomic.LoadInt32(&wr.retry_attempts) > wr.max_retry_attempts {
				return 0, fmt.Errorf("Exceeded max retry attempts")
			}
		}

		// 409 error

		if update_rsp.StatusCode == 409 {
			return wr.Write(ctx, uri, fh)
		}

		// rate limit

		rate := ratelimit_err.Rate
		reset := rate.Reset
		then := reset.Unix()

		now := time.Now()
		ts := now.Unix()

		wait := then - ts
		duration := time.Duration(time.Duration(wait) * time.Second)

		time.Sleep(duration)

		return wr.Write(ctx, uri, fh)
	}

	return 0, nil

}

func (wr *GitHubAPIWriter) WriterURI(ctx context.Context, key string) string {

	uri := key

	if wr.prefix != "" {
		uri = filepath.Join(wr.prefix, key)
	}

	return uri
}

func (wr *GitHubAPIWriter) Flush(ctx context.Context) error {
	return nil
}

func (wr *GitHubAPIWriter) Close(ctx context.Context) error {
	return nil
}

func (wr *GitHubAPIWriter) SetLogger(ctx context.Context, logger *log.Logger) error {
	return nil
}
