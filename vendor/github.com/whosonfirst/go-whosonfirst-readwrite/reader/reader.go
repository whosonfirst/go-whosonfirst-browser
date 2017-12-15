package reader

import (
	"errors"
	"io"
	"net/url"
)

type Reader interface {
	Read(string) (io.ReadCloser, error)
}

func Sources() []string {
     return []string{ "fs", "http", "s3" }
}

func NewReaderFromSource(source string, args ...interface{}) (Reader, error) {

	var r Reader
	var err error

	switch source {
	case "fs":

		if len(args) == 0 {
			err = errors.New("Missing FS root")
		}

		root := args[0].(string)

		r, err = NewFSReader(root)

	case "http":

		if len(args) == 0 {
			err = errors.New("Missing HTTP root")
		}

		root, err := url.Parse(args[0].(string))

		if err != nil {
			return nil, err
		}

		r, err = NewHTTPReader(root)

	case "s3":

		if len(args) < 4 {
			err = errors.New("Insufficient S3 arguments")
		}

		bucket := args[0].(string)
		prefix := args[1].(string)
		region := args[2].(string)
		creds := args[3].(string)

		cfg := S3Config{
			Bucket:      bucket,
			Prefix:      prefix,
			Region:      region,
			Credentials: creds,
		}

		r, err = NewS3Reader(cfg)
	default:
		err = errors.New("Unknown or invalid source")
	}

	return r, err
}
