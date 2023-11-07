package writer

// This is basically a thin wrapper on top of this:
// https://github.com/google/go-github/blob/master/example/commitpr/main.go

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
	"sync"
	"time"

	"github.com/google/go-github/v48/github"
	wof_writer "github.com/whosonfirst/go-writer/v3"
	"golang.org/x/oauth2"
)

const GITHUBAPI_PR_SCHEME string = "githubapi-pr"

// base_ is the thing a PR is being created "against"
// pr_ is the thing where the PR is being created

type GitHubAPIPullRequestWriter struct {
	wof_writer.Writer
	base_owner     string
	base_repo      string
	base_branch    string
	pr_owner       string
	pr_repo        string
	pr_branch      string
	pr_author      string
	pr_email       string
	pr_title       string
	pr_description string
	pr_entries     []*github.TreeEntry
	pr_ensure_repo bool
	prefix         string
	client         *github.Client
	user           *github.User
	logger         *log.Logger
	mutex          *sync.RWMutex
}

func init() {

	ctx := context.Background()
	wof_writer.RegisterWriter(ctx, GITHUBAPI_PR_SCHEME, NewGitHubAPIPullRequestWriter)
}

func NewGitHubAPIPullRequestWriter(ctx context.Context, uri string) (wof_writer.Writer, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	base_owner := u.Host

	path := strings.TrimLeft(u.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) != 1 {
		return nil, errors.New("Invalid path")
	}

	base_repo := parts[0]
	base_branch := DEFAULT_BRANCH

	q := u.Query()

	token := q.Get("access_token")

	prefix := q.Get("prefix")
	branch := q.Get("branch")

	if token == "" {
		return nil, errors.New("Missing access token")
	}

	if branch != "" {
		base_branch = branch
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	users := client.Users
	user, _, err := users.Get(ctx, "")

	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve user for token, %w", err)
	}
	pr_owner := q.Get("pr-owner")

	if pr_owner == "" {
		pr_owner = base_owner
	}

	pr_repo := q.Get("pr-repo")

	if pr_repo == "" {
		pr_repo = base_repo
	}

	pr_branch := q.Get("pr-branch")

	if pr_branch == "" {
		return nil, fmt.Errorf("Invalid pr-branch")
	}

	if pr_branch == branch {
		return nil, fmt.Errorf("pr-branch can not be the same as branch")
	}

	pr_title := q.Get("pr-title")

	if pr_title == "" {
		return nil, fmt.Errorf("Invalid pr-title")
	}

	pr_description := q.Get("pr-description")

	if pr_title == "" {
		return nil, fmt.Errorf("Invalid pr-title")
	}

	pr_author := q.Get("pr-author")

	if pr_author == "" {
		pr_author = user.GetName()
	}

	if pr_author == "" {
		return nil, fmt.Errorf("Invalid pr-author argument")
	}

	pr_email := q.Get("pr-email")

	if pr_email == "" {
		pr_email = user.GetEmail()
	}

	if pr_email == "" {
		return nil, fmt.Errorf("Invalid pr-email argument")
	}

	str_ensure_repo := q.Get("ensure-repo")

	pr_ensure_repo := false

	if str_ensure_repo != "" {

		ensure_repo, err := strconv.ParseBool(str_ensure_repo)

		if err != nil {
			return nil, fmt.Errorf("Invalid ensure-repo argument, %w", err)
		}

		pr_ensure_repo = ensure_repo
	}

	pr_entries := []*github.TreeEntry{}

	mutex := new(sync.RWMutex)
	logger := log.Default()

	wr := &GitHubAPIPullRequestWriter{
		client:         client,
		user:           user,
		base_owner:     base_owner,
		base_repo:      base_repo,
		base_branch:    base_branch,
		pr_owner:       pr_owner,
		pr_repo:        pr_repo,
		pr_branch:      pr_branch,
		pr_author:      pr_author,
		pr_email:       pr_email,
		pr_title:       pr_title,
		pr_description: pr_description,
		pr_ensure_repo: pr_ensure_repo,
		pr_entries:     pr_entries,
		prefix:         prefix,
		logger:         logger,
		mutex:          mutex,
	}

	return wr, nil
}

func (wr *GitHubAPIPullRequestWriter) Write(ctx context.Context, uri string, r io.ReadSeeker) (int64, error) {

	// Something something something account for cases with a bazillion commits and not keeping
	// everything in memory until we call Close(). One option would be to keep a local map of io.ReadSeeker
	// instances but then we will just have filehandle exhaustion problems. Add option to write to
	// disk or something like a SQLite database (allowing a custom DSN to determine whether to write to
	// disk or memory) ?

	body, err := io.ReadAll(r)

	if err != nil {
		return 0, err
	}

	wr_uri := wr.WriterURI(ctx, uri)

	e := &github.TreeEntry{
		Path:    github.String(wr_uri),
		Type:    github.String("blob"),
		Content: github.String(string(body)),
		Mode:    github.String("100644"),
	}

	wr.mutex.Lock()
	defer wr.mutex.Unlock()

	wr.pr_entries = append(wr.pr_entries, e)

	return 0, nil
}

func (wr *GitHubAPIPullRequestWriter) Flush(ctx context.Context) error {

	wr.mutex.Lock()
	defer wr.mutex.Unlock()

	if len(wr.pr_entries) == 0 {
		return nil
	}

	if wr.pr_ensure_repo {

		err := wr.ensureRepo(ctx)

		if err != nil {
			return fmt.Errorf("Failed to ensure repo, %w", err)
		}
	}

	ref, err := wr.getRef(ctx)

	if err != nil {

		if err != nil {
			return fmt.Errorf("Failed to get ref, %w", err)
		}
	}

	tree, _, err := wr.client.Git.CreateTree(ctx, wr.pr_owner, wr.pr_repo, *ref.Object.SHA, wr.pr_entries)

	if err != nil {
		return fmt.Errorf("Failed to create tree, %w", err)
	}

	err = wr.pushCommit(ctx, ref, tree)

	if err != nil {
		return fmt.Errorf("Failed to push commit, %w", err)
	}

	return nil
}

func (wr *GitHubAPIPullRequestWriter) Close(ctx context.Context) error {

	err := wr.Flush(ctx)

	if err != nil {
		return fmt.Errorf("Failed to flush records, %w", err)
	}

	err = wr.createPR(ctx)

	if err != nil {
		return fmt.Errorf("Failed to create PR, %w", err)
	}

	return nil
}

func (wr *GitHubAPIPullRequestWriter) SetLogger(ctx context.Context, logger *log.Logger) error {
	wr.logger = logger
	return nil
}

func (wr *GitHubAPIPullRequestWriter) WriterURI(ctx context.Context, key string) string {

	uri := key

	if wr.prefix != "" {
		uri = filepath.Join(wr.prefix, key)
	}

	return uri
}

func (wr *GitHubAPIPullRequestWriter) ensureRepo(ctx context.Context) error {

	// https://github.com/google/go-github/blob/master/github/repos_forks.go

	_, _, err := wr.client.Repositories.Get(ctx, wr.pr_owner, wr.pr_repo)

	if err == nil {
		return nil
	}

	if wr.base_owner == wr.pr_owner {
		return fmt.Errorf("Can not fork %s from %s to %s", wr.base_repo, wr.base_owner, wr.pr_owner)
	}

	fork_opts := &github.RepositoryCreateForkOptions{
		Organization: wr.pr_owner,
	}

	repo, _, err := wr.client.Repositories.CreateFork(ctx, wr.base_owner, wr.base_repo, fork_opts)

	if err != nil {

		_, accepted_err := err.(*github.AcceptedError)

		if !accepted_err {
			return fmt.Errorf("Failed to create fork, %w", err)
		}

		time.Sleep(1)
	}

	if wr.base_repo != wr.pr_repo {

		repo.Name = &wr.pr_repo

		_, _, err := wr.client.Repositories.Edit(ctx, wr.pr_owner, wr.base_repo, repo)

		if err != nil {
			return fmt.Errorf("Failed to rename fork, %w", err)
		}
	}

	return nil
}

func (wr *GitHubAPIPullRequestWriter) getRef(ctx context.Context) (*github.Reference, error) {

	base_branch := fmt.Sprintf("refs/heads/%s", wr.base_branch)
	pr_branch := fmt.Sprintf("refs/heads/%s", wr.pr_branch)

	pr_ref, _, _ := wr.client.Git.GetRef(ctx, wr.pr_owner, wr.pr_repo, pr_branch)

	if pr_ref != nil {
		return pr_ref, nil
	}

	base_ref, _, err := wr.client.Git.GetRef(ctx, wr.pr_owner, wr.pr_repo, base_branch)

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve base branch '%s' for %s/%s, %w", base_branch, wr.pr_owner, wr.pr_repo, err)
	}

	new_ref := &github.Reference{Ref: github.String(pr_branch), Object: &github.GitObject{SHA: base_ref.Object.SHA}}

	pr_ref, _, err = wr.client.Git.CreateRef(ctx, wr.pr_owner, wr.pr_repo, new_ref)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new ref, %w", err)
	}

	return pr_ref, err
}

// pushCommit creates the commit in the given reference using the given tree.
func (wr *GitHubAPIPullRequestWriter) pushCommit(ctx context.Context, ref *github.Reference, tree *github.Tree) error {

	// Get the parent commit to attach the commit to.

	list_opts := &github.ListOptions{}

	parent, _, err := wr.client.Repositories.GetCommit(ctx, wr.pr_owner, wr.pr_repo, *ref.Object.SHA, list_opts)

	if err != nil {
		return fmt.Errorf("Failed to determine parent commit, %w", err)
	}

	// This is not always populated, but is needed.
	parent.Commit.SHA = parent.SHA

	// Create the commit using the tree.
	date := time.Now()

	author := &github.CommitAuthor{
		Date:  &date,
		Name:  &wr.pr_author,
		Email: &wr.pr_email,
	}

	parents := []*github.Commit{
		parent.Commit,
	}

	commit := &github.Commit{
		Author:  author,
		Message: &wr.pr_description,
		Tree:    tree,
		Parents: parents,
	}

	newCommit, _, err := wr.client.Git.CreateCommit(ctx, wr.pr_owner, wr.pr_repo, commit)

	if err != nil {
		return fmt.Errorf("Failed to create commit, %w", err)
	}

	// Attach the commit to the main branch.
	ref.Object.SHA = newCommit.SHA

	_, _, err = wr.client.Git.UpdateRef(ctx, wr.pr_owner, wr.pr_repo, ref, false)

	if err != nil {
		return fmt.Errorf("Failed to update ref, %w", err)
	}

	return nil
}

// createPR creates a pull request. Based on: https://godoc.org/github.com/google/go-github/github#example-PullRequestsService-Create
func (wr *GitHubAPIPullRequestWriter) createPR(ctx context.Context) error {

	head := wr.pr_branch

	if wr.pr_owner != wr.base_owner {
		head = fmt.Sprintf("%s:%s", wr.pr_owner, wr.pr_branch)
	}

	new_pr := &github.NewPullRequest{
		Title:               &wr.pr_title,
		Head:                &head, // &wr.pr_branch,
		Base:                &wr.base_branch,
		Body:                &wr.pr_description,
		MaintainerCanModify: github.Bool(true),
	}

	// Maybe capture (new) pr here and store in ctx if an application wants to retrieve it later?

	_, _, err := wr.client.PullRequests.Create(ctx, wr.base_owner, wr.base_repo, new_pr)

	if err != nil {
		return fmt.Errorf("Failed to create pull request, %w", err)
	}

	return nil
}
