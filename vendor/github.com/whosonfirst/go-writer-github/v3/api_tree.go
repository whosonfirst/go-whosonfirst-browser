package writer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v48/github"
	wof_writer "github.com/whosonfirst/go-writer/v3"
	"golang.org/x/oauth2"
)

const GITHUBAPI_TREE_SCHEME string = "githubapi-tree"

type GitHubAPITreeWriter struct {
	wof_writer.Writer
	base_owner         string
	base_repo          string
	base_branch        string
	commit_owner       string
	commit_repo        string
	commit_branch      string
	commit_author      string
	commit_email       string
	commit_description string
	commit_entries     []*github.TreeEntry
	commit_ensure_repo bool
	prefix             string
	client             *github.Client
	user               *github.User
	logger             *log.Logger
	mutex              *sync.RWMutex
}

func init() {

	ctx := context.Background()
	wof_writer.RegisterWriter(ctx, GITHUBAPI_TREE_SCHEME, NewGitHubAPITreeWriter)
}

func NewGitHubAPITreeWriter(ctx context.Context, uri string) (wof_writer.Writer, error) {

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

	commit_owner := base_owner
	commit_repo := base_repo
	commit_branch := base_branch

	to_branch := q.Get("to-branch")

	if to_branch != "" {
		commit_branch = to_branch
	}

	commit_description := q.Get("description")

	commit_author := q.Get("author")

	if commit_author == "" {
		commit_author = user.GetName()
	}

	if commit_author == "" {
		return nil, fmt.Errorf("Invalid author")
	}

	commit_email := q.Get("email")

	if commit_email == "" {
		commit_email = user.GetEmail()
	}

	if commit_email == "" {
		return nil, fmt.Errorf("Invalid email address")
	}

	commit_entries := []*github.TreeEntry{}

	logger := log.Default()

	mutex := new(sync.RWMutex)

	wr := &GitHubAPITreeWriter{
		client:             client,
		user:               user,
		base_owner:         base_owner,
		base_repo:          base_repo,
		base_branch:        base_branch,
		commit_owner:       commit_owner,
		commit_repo:        commit_repo,
		commit_branch:      commit_branch,
		commit_author:      commit_author,
		commit_email:       commit_email,
		commit_description: commit_description,
		commit_entries:     commit_entries,
		prefix:             prefix,
		logger:             logger,
		mutex:              mutex,
	}

	return wr, nil
}

func (wr *GitHubAPITreeWriter) Write(ctx context.Context, uri string, r io.ReadSeeker) (int64, error) {

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

	wr.commit_entries = append(wr.commit_entries, e)

	wr.logger.Printf("Add %s\n", wr_uri)
	return 0, nil
}

func (wr *GitHubAPITreeWriter) Flush(ctx context.Context) error {

	wr.mutex.Lock()
	defer wr.mutex.Unlock()

	if len(wr.commit_entries) == 0 {
		return nil
	}

	ref, err := wr.getRef(ctx)

	if err != nil {

		if err != nil {
			return fmt.Errorf("Failed to get ref, %w", err)
		}
	}

	tree, _, err := wr.client.Git.CreateTree(ctx, wr.commit_owner, wr.commit_repo, *ref.Object.SHA, wr.commit_entries)

	if err != nil {
		return fmt.Errorf("Failed to create tree, %w", err)
	}

	err = wr.pushCommit(ctx, ref, tree)

	if err != nil {
		return fmt.Errorf("Failed to push commit, %w", err)
	}

	wr.commit_entries = []*github.TreeEntry{}
	return nil
}

func (wr *GitHubAPITreeWriter) Close(ctx context.Context) error {
	return wr.Flush(ctx)
}

func (wr *GitHubAPITreeWriter) SetLogger(ctx context.Context, logger *log.Logger) error {
	wr.logger = logger
	return nil
}

func (wr *GitHubAPITreeWriter) WriterURI(ctx context.Context, key string) string {

	uri := key

	if wr.prefix != "" {
		uri = filepath.Join(wr.prefix, key)
	}

	return uri
}

func (wr *GitHubAPITreeWriter) getRef(ctx context.Context) (*github.Reference, error) {

	base_branch := fmt.Sprintf("refs/heads/%s", wr.base_branch)
	commit_branch := fmt.Sprintf("refs/heads/%s", wr.commit_branch)

	commit_ref, _, _ := wr.client.Git.GetRef(ctx, wr.commit_owner, wr.commit_repo, commit_branch)

	if commit_ref != nil {
		return commit_ref, nil
	}

	base_ref, _, err := wr.client.Git.GetRef(ctx, wr.commit_owner, wr.commit_repo, base_branch)

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve base branch '%s' for %s/%s, %w", base_branch, wr.commit_owner, wr.commit_repo, err)
	}

	new_ref := &github.Reference{Ref: github.String(commit_branch), Object: &github.GitObject{SHA: base_ref.Object.SHA}}

	commit_ref, _, err = wr.client.Git.CreateRef(ctx, wr.commit_owner, wr.commit_repo, new_ref)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new ref, %w", err)
	}

	return commit_ref, err
}

// pushCommit creates the commit in the given reference using the given tree.
func (wr *GitHubAPITreeWriter) pushCommit(ctx context.Context, ref *github.Reference, tree *github.Tree) error {

	// Get the parent commit to attach the commit to.

	list_opts := &github.ListOptions{}

	parent, _, err := wr.client.Repositories.GetCommit(ctx, wr.commit_owner, wr.commit_repo, *ref.Object.SHA, list_opts)

	if err != nil {
		return fmt.Errorf("Failed to determine parent commit, %w", err)
	}

	// This is not always populated, but is needed.
	parent.Commit.SHA = parent.SHA

	// Create the commit using the tree.
	date := time.Now()

	author := &github.CommitAuthor{
		Date:  &date,
		Name:  &wr.commit_author,
		Email: &wr.commit_email,
	}

	parents := []*github.Commit{
		parent.Commit,
	}

	commit := &github.Commit{
		Author:  author,
		Message: &wr.commit_description,
		Tree:    tree,
		Parents: parents,
	}

	newCommit, _, err := wr.client.Git.CreateCommit(ctx, wr.commit_owner, wr.commit_repo, commit)

	if err != nil {
		return fmt.Errorf("Failed to create commit, %w", err)
	}

	// Attach the commit to the main branch.
	ref.Object.SHA = newCommit.SHA

	_, _, err = wr.client.Git.UpdateRef(ctx, wr.commit_owner, wr.commit_repo, ref, false)

	if err != nil {
		return fmt.Errorf("Failed to update ref, %w", err)
	}

	return nil
}
