package writer

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	wof_writer "github.com/whosonfirst/go-writer"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	_ "log"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

type GitHubAPIWriterCommitTemplates struct {
	New    string
	Update string
}

type GitHubAPIWriter struct {
	wof_writer.Writer
	owner     string
	repo      string
	branch    string
	prefix    string
	client    *github.Client
	user      *github.User
	throttle  <-chan time.Time
	templates *GitHubAPIWriterCommitTemplates
}

func init() {

	ctx := context.Background()
	wof_writer.RegisterWriter(ctx, "githubapi", NewGitHubAPIWriter)
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
	branch := "master"

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

	rate := time.Second / 3
	throttle := time.Tick(rate)

	wr := &GitHubAPIWriter{
		client:    client,
		owner:     u.Host,
		user:      user,
		repo:      repo,
		branch:    branch,
		prefix:    prefix,
		templates: templates,
		throttle:  throttle,
	}

	return wr, nil
}

func (wr *GitHubAPIWriter) Write(ctx context.Context, uri string, fh io.ReadCloser) error {

	<-wr.throttle

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return err
	}

	url := wr.URI(uri)

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

	get_opts := &github.RepositoryContentGetOptions{}

	get_rsp, _, _, err := wr.client.Repositories.GetContents(ctx, wr.owner, wr.repo, url, get_opts)

	if err == nil {
		commit_msg = fmt.Sprintf(wr.templates.Update, url)
		update_opts.Message = github.String(commit_msg)
		update_opts.SHA = get_rsp.SHA
	}

	_, _, err = wr.client.Repositories.UpdateFile(ctx, wr.owner, wr.repo, url, update_opts)

	if err != nil {
		return err
	}

	return nil

}

func (wr *GitHubAPIWriter) URI(key string) string {

	uri := key

	if wr.prefix != "" {
		uri = filepath.Join(wr.prefix, key)
	}

	return uri
}
