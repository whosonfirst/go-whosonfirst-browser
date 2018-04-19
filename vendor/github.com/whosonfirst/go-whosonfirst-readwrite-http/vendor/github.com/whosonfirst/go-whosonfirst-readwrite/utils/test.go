package utils

import (
	"github.com/whosonfirst/go-whosonfirst-readwrite/reader"
)

func TestReader(r reader.Reader, key string) (bool, error) {

	_, err := r.Read(key)

	if err != nil {
		return false, err
	}

	return true, nil
}
