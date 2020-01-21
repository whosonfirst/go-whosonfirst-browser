# go-writer-github

GitHub API support for the [go-writer](https://github.com/whosonfirst/go-writer) `Writer` interface.

This package supersedes the [go-whosonfirst-readwrite-github](https://github.com/whosonfirst/go-whosonfirst-readwrite-github) package.

## Example

```
package main

import (
	"context"
	"flag"
	"fmt"	
	"github.com/whosonfirst/go-writer"
	_ "github.com/whosonfirst/go-writer-github"		
	"io/ioutil"
	"strings"
	"time"
)

func main() {

	owner := "example"
	repo := "example-repo"
	access_token := "s33kret"
	branch := "example"

	ctx := context.Background()
	
	source := fmt.Sprintf("githubapi://%s/%s?access_token=%s&branch=%s", owner, repo, access_token, branch)
			
	wr, _ := writer.NewWriter(ctx, source)

	now := time.Now()
	msg := fmt.Sprintf("This is a test: %v", now)
	
	br := strings.NewReader(msg)
	fh := ioutil.NopCloser(br)

	wr.Write(ctx, "TEST.md", fh)
}
```

_Error handling has been omitted for brevity._

## URIs

The URI structure for `go-writer-github` (API) writers is:

| Component | Value | Notes |
| --- | --- | --- |
| Scheme | `githubapi` | |
| Host | string | A valid GitHub owner or organization |
| Path | string | A valid GitHub respository |
| Query | | _described below_ |

The following query parameters are supported for `go-writer-github` writers:

| Name | Value | Required | Notes |
| --- | --- | --- | --- |
| access_token | string | yes | A valid GitHub API access token |
| branch | string | no | A valid Git repository branch. Default is `master`. |
| new | string | no | A valid string for formatting new file commit messages. Default is `Created %s`. |
| update | string | no | A valid string for formatting updated file commit messages. Default is `Updated %s`. |

## See also

* https://github.com/whosonfirst/go-writer
* https://github.com/google/go-github