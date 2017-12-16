package reader

// maybe also make a GH API reader...
// https://developer.github.com/v3/repos/contents/#get-contents

import (
	"fmt"
	"io"
	_ "log"
	"net/http"
)

type GitHubReader struct {
	Reader
	repo string
}

func NewGitHubReader(repo string) (Reader, error) {

	r := GitHubReader{
		repo: repo,
	}

	return &r, nil
}

func (r *GitHubReader) Read(key string) (io.ReadCloser, error) {

	url := fmt.Sprintf("https://raw.githubusercontent.com/whosonfirst-data/%s/%s", r.repo, key)

	rsp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	return rsp.Body, nil
}
