package dynamodb

import (
	"context"
	"fmt"
	"github.com/aaronland/go-aws-session"
	"github.com/aws/aws-sdk-go/aws"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	aws_dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"net/url"
)

// 'awsdynamodb://findinaid?region=us-west-2&endpoint=http://localhost:8000&credentials=static:local:local:local'

func NewClientWithURI(ctx context.Context, uri string) (*aws_dynamodb.DynamoDB, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %v", err)
	}

	// table_name := u.Host

	q := u.Query()

	// partition_key := q.Get("partition_key")

	region := q.Get("region")
	endpoint := q.Get("endpoint")

	credentials := q.Get("credentials")

	cfg, err := session.NewConfigWithCredentialsAndRegion(credentials, region)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new session for credentials '%s', %w", credentials, err)
	}

	if endpoint != "" {
		cfg.Endpoint = aws.String(endpoint)
	}

	sess, err := aws_session.NewSession(cfg)

	if err != nil {
		return nil, fmt.Errorf("Failed to create AWS session, %w", err)
	}

	client := aws_dynamodb.New(sess)
	return client, nil
}
