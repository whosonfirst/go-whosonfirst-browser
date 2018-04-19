package config

import (
       "errors"
       "strings"
)

type S3Config struct {
	Bucket      string
	Prefix      string
	Region      string
	Credentials string // see notes below
}

func NewS3ConfigFromString(cfg string) (*S3Config, error) {

	parts := strings.Split(cfg, " ")

	if len(parts) != 4 {
		return nil, errors.New("Invalid count for config")
	}

	s3_config := S3Config{
		Bucket:      "",
		Prefix:      "",
		Region:      "",
		Credentials: "",
	}

	for _, p := range parts {

		p = strings.Trim(p, " ")
		kv := strings.Split(p, "=")

		if len(kv) != 2 {
			return nil, errors.New("Invalid count for config block")
		}

		switch kv[0] {
		case "bucket":
			s3_config.Bucket = kv[1]
		case "prefix":
			s3_config.Prefix = kv[1]
		case "region":
			s3_config.Region = kv[1]
		case "credentials":
			s3_config.Credentials = kv[1]
		default:
			return nil, errors.New("Invalid key for config block")
		}
	}

	return &s3_config, nil
}
