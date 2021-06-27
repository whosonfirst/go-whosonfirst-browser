package cache

import (
	"bytes"
	"context"
	"fmt"
	"io"
	_ "log"
)

type testCacheOptions struct {
	AllowCacheMiss bool
}

func testCache(ctx context.Context, c Cache, opts *testCacheOptions) error {

	k := "test"
	v := "test"

	fh1, err := ReadSeekCloserFromString(v)

	if err != nil {
		return fmt.Errorf("Failed to create ReadSeekCloser, %v", err)
	}

	fh2, err := c.Set(ctx, k, fh1)

	if err != nil {
		return fmt.Errorf("Failed to set %s (%s), %v", k, v, err)
	}

	equals, err := compareReadSeekClosers(fh1, fh2)

	if err != nil {
		return fmt.Errorf("Failed to compare filehandles, %v", err)
	}

	if !equals {
		return fmt.Errorf("Filehandles 1 and 2 failed equality test")
	}

	fh3, err := c.Get(ctx, k)

	if err != nil && !IsCacheMiss(err) {
		return fmt.Errorf("Failed to get cache value for '%s', %v", k, err)
	}

	if err != nil {

		if !opts.AllowCacheMiss {
			return fmt.Errorf("Unexpected cache miss for '%s'", k)
		}

	} else {

		equals, err = compareReadSeekClosers(fh1, fh3)

		if err != nil {
			return fmt.Errorf("Failed to compare filehandles 1 and 3, %v", err)
		}

		if !equals {
			return fmt.Errorf("Filehandles 1 and 3 failed equality test")
		}
	}

	return nil
}

func compareReadSeekClosers(fh1 io.ReadSeekCloser, fh2 io.ReadSeekCloser) (bool, error) {

	fh1.Seek(0, 0)
	fh2.Seek(0, 0)

	defer func() {
		fh1.Seek(0, 0)
		fh2.Seek(0, 0)
	}()

	b1, err := io.ReadAll(fh1)

	if err != nil {
		return false, fmt.Errorf("Failed to read fh1, %v", err)
	}

	fh1.Seek(0, 0)

	b2, err := io.ReadAll(fh2)

	if err != nil {
		return false, fmt.Errorf("Failed to read fh2, %v", err)
	}

	fh2.Seek(0, 0)

	// log.Printf("1 '%s'\n", string(b1))
	// log.Printf("2 '%s'\n", string(b2))

	equals := bytes.Compare(b1, b2)

	if equals != 0 {
		return false, nil
	}

	return true, nil
}
