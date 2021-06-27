package repo

import (
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"path/filepath"
	"strings"
)

func cacheKeyFromURI(str_uri string) (string, error) {

	id, uri_args, err := uri.ParseURI(str_uri)

	if err != nil {
		return "", err
	}

	rel_path, err := uri.Id2RelPath(id, uri_args)

	if err != nil {
		return "", err
	}

	return cacheKeyFromRelPath(rel_path)
}

func cacheKeyFromRelPath(rel_path string) (string, error) {

	ext := filepath.Ext(rel_path)
	rel_path = strings.Replace(rel_path, ext, "", 1)

	key := fmt.Sprintf("%s.json", rel_path)
	return key, nil
}
