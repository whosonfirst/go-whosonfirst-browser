package writer

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v48/github"
	wof_writer "github.com/whosonfirst/go-writer/v3"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const GITHUBAPI_BRANCH_SCHEME string = "githubapi-branch"

type GitHubAPIBranchWriter struct {
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
	merge_branch       bool
	remove_branch      bool
	prefix             string
	client             *github.Client
	user               *github.User
	logger             *log.Logger
	mutex              *sync.RWMutex
}

func init() {
	ctx := context.Background()
	wof_writer.RegisterWriter(ctx, GITHUBAPI_BRANCH_SCHEME, NewGitHubAPIBranchWriter)
}

func NewGitHubAPIBranchWriter(ctx context.Context, uri string) (wof_writer.Writer, error) {

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

	to_branch := q.Get("to-branch")

	if to_branch == "" {
		return nil, fmt.Errorf("Invalid to-branch")
	}

	if to_branch == base_branch {
		return nil, fmt.Errorf("Commit branch can not be the same as base branch")
	}

	commit_branch := to_branch

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

	merge_branch := false
	remove_branch := false

	str_merge := q.Get("merge")
	str_remove := q.Get("remove-on-merge")

	if str_merge != "" {

		merge, err := strconv.ParseBool(str_merge)

		if err != nil {
			return nil, fmt.Errorf("Invalid merge parameter, %w", err)
		}

		merge_branch = merge
	}

	if str_remove != "" {

		remove, err := strconv.ParseBool(str_remove)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse '%s', %v", str_remove, err)
		}

		remove_branch = remove
	}

	logger := log.New(io.Discard, "", 0)

	mutex := new(sync.RWMutex)

	wr := &GitHubAPIBranchWriter{
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
		merge_branch:       merge_branch,
		remove_branch:      remove_branch,
	}

	return wr, nil
}

func (wr *GitHubAPIBranchWriter) Write(ctx context.Context, uri string, r io.ReadSeeker) (int64, error) {

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

	wr.logger.Printf("Add %s/%s @%s\n", wr.base_repo, wr_uri, wr.commit_branch)
	return 0, nil
}

func (wr *GitHubAPIBranchWriter) Flush(ctx context.Context) error {

	wr.logger.Printf("Flush %d entries for %s %s\n", len(wr.commit_entries), wr.commit_repo, wr.commit_branch)

	wr.mutex.Lock()
	defer wr.mutex.Unlock()

	if len(wr.commit_entries) == 0 {
		wr.logger.Printf("No entries to flush for %s @%s\n", wr.commit_repo, wr.commit_branch)
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
		return fmt.Errorf("Failed to create branch, %w", err)
	}

	err = wr.pushCommit(ctx, ref, tree)

	if err != nil {
		return fmt.Errorf("Failed to push commit, %w", err)
	}

	wr.commit_entries = []*github.TreeEntry{}
	return nil
}

func (wr *GitHubAPIBranchWriter) Close(ctx context.Context) error {

	err := wr.Flush(ctx)

	if err != nil {
		return fmt.Errorf("Failed to flush writer, %w", err)
	}

	if !wr.merge_branch {
		return nil
	}

	err = wr.mergeBranch(ctx)

	if err != nil {
		return fmt.Errorf("Failed to merge branch, %w", err)
	}

	if !wr.remove_branch {
		return nil
	}

	err = wr.removeBranch(ctx)

	if err != nil {
		return fmt.Errorf("Failed to remove branch, %w", err)
	}

	return nil
}

func (wr *GitHubAPIBranchWriter) SetLogger(ctx context.Context, logger *log.Logger) error {
	wr.logger = logger
	return nil
}

func (wr *GitHubAPIBranchWriter) WriterURI(ctx context.Context, key string) string {

	uri := key

	if wr.prefix != "" {
		uri = filepath.Join(wr.prefix, key)
	}

	return uri
}

func (wr *GitHubAPIBranchWriter) getRef(ctx context.Context) (*github.Reference, error) {

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

	new_ref := &github.Reference{
		Ref: github.String(commit_branch),
		Object: &github.GitObject{
			SHA: base_ref.Object.SHA,
		},
	}

	commit_ref, _, err = wr.client.Git.CreateRef(ctx, wr.commit_owner, wr.commit_repo, new_ref)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new ref, %w", err)
	}

	return commit_ref, err
}

// pushCommit creates the commit in the given reference using the given branch.
func (wr *GitHubAPIBranchWriter) pushCommit(ctx context.Context, ref *github.Reference, tree *github.Tree) error {

	wr.logger.Printf("Push commit for %s @%s\n", wr.commit_repo, wr.commit_branch)

	// Get the parent commit to attach the commit to.

	list_opts := &github.ListOptions{}

	parent, _, err := wr.client.Repositories.GetCommit(ctx, wr.commit_owner, wr.commit_repo, *ref.Object.SHA, list_opts)

	if err != nil {
		return fmt.Errorf("Failed to determine parent commit, %w", err)
	}

	// This is not always populated, but is needed.
	parent.Commit.SHA = parent.SHA

	// Create the commit using the branch.
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

func (wr *GitHubAPIBranchWriter) mergeBranch(ctx context.Context) error {

	commit_msg := fmt.Sprintf("Merge %s", wr.commit_branch)

	wr.logger.Println(commit_msg)

	req := &github.RepositoryMergeRequest{
		Base:          &wr.base_branch,
		Head:          &wr.commit_branch,
		CommitMessage: &commit_msg,
	}

	_, _, err := wr.client.Repositories.Merge(ctx, wr.base_owner, wr.base_repo, req)

	if err != nil {
		return fmt.Errorf("Failed to merge branch, %w", err)
	}

	return nil
}

func (wr *GitHubAPIBranchWriter) removeBranch(ctx context.Context) error {

	ref := fmt.Sprintf("heads/%s", wr.commit_branch)

	wr.logger.Printf("Remove %s\n", ref)

	_, err := wr.client.Git.DeleteRef(ctx, wr.base_owner, wr.base_repo, ref)

	if err != nil {
		return fmt.Errorf("Failed to remove branch, %w", err)
	}

	return nil
}
