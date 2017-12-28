package reader

// maybe also make a GH API reader...
// https://developer.github.com/v3/repos/contents/#get-contents

import (
	"errors"
	"fmt"
	"io"
	_ "log"
	"net/http"
)

type GitHubReader struct {
	Reader
	repo   string
	branch string
}

func NewGitHubReader(repo string, branch string) (Reader, error) {

	r := GitHubReader{
		repo:   repo,
		branch: branch,
	}

	return &r, nil
}

func (r *GitHubReader) Read(key string) (io.ReadCloser, error) {

	url := fmt.Sprintf("https://raw.githubusercontent.com/whosonfirst-data/%s/%s/data/%s", r.repo, r.branch, key)

	// log.Println("READ", key, url)

	rsp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != 200 {
		return nil, errors.New(rsp.Status)
	}

	return rsp.Body, nil
}
