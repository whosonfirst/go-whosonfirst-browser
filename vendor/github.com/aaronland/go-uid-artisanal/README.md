# go-uid-artisanal

Work in progress.

## Example

```
package main

import (
	"context"
	_ "github.com/aaronland/go-brooklynintegers-api"
	"github.com/aaronland/go-uid"
	"log"
)

func main() {

	ctx := context.Background()

	opts := &ArtisanalProviderURIOptions{
		Pool:    "memory://",
		Minimum: 5,
		Clients: []string{
			"brooklynintegers://",
		},
	}

	uri, _ := NewArtisanalProviderURI(opts)
	str_uri := uri.String()

	// str_uri ends up looking like this:
	// artisanal:?client=brooklynintegers%3A%2F%2F&minimum=5&pool=memory%3A%2F%2F
	
	pr, _ := uid.NewProvider(ctx, str_uri)
	id, _ := pr.UID()

	log.Println(id.String())
}

```

## See also

* https://github.com/aaronland/go-uid
* https://github.com/aaronland/go-pool
* https://github.com/aaronland/go-artisanal-integers
* https://github.com/aaronland/go-artisanal-integers-proxy/service