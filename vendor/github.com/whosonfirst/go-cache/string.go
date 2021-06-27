package cache

import (
	"context"
	"github.com/whosonfirst/go-ioutil"
	"io"
	"strings"
)

func ReadSeekCloserFromString(v string) (io.ReadSeekCloser, error) {
	sr := strings.NewReader(v)
	return ioutil.NewReadSeekCloser(sr)
}

func SetString(ctx context.Context, c Cache, k string, v string) (string, error) {

	fh, err := ReadSeekCloserFromString(v)

	if err != nil {
		return "", err
	}

	fh, err = c.Set(ctx, k, fh)

	if err != nil {
		return "", err
	}

	defer fh.Close()

	return toString(fh)
}

func GetString(ctx context.Context, c Cache, k string) (string, error) {

	fh, err := c.Get(ctx, k)

	if err != nil {
		return "", err
	}

	defer fh.Close()

	return toString(fh)
}

func toString(fh io.Reader) (string, error) {

	b, err := io.ReadAll(fh)

	if err != nil {
		return "", err
	}

	return string(b), nil
}
