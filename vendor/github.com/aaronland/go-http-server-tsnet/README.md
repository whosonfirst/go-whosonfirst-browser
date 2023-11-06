# go-http-server-tsnet

Experimental Tailscale `tsnet` implementation of the `aaronland/go-http-server` interfaces.

This package creates a Tailscale virtual private service for HTTP resources. It wraps all the `http.Handler` instances in a middleware wrapper that derives and stores the Tailscale user and machine that is accessing the server.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/aaronland/go-http-server-tsnet.svg)](https://pkg.go.dev/github.com/aaronland/go-http-server-tsnet)

## Known-knowns

* Automagic TLS (port 443) certificates don't work with errors whose words I understand but not their meaning.
* Under the hood there is extra undocumented hoop-jumping to make passing TS auth keys as query parameters (to the server URI constructor) work.

## Example

_Error handling removed for the sake of brevity._

```
$> go run cmd/example/main.go \
	-server-uri 'tsnet://testing:80?auth-key={TS_AUTH_KEY}'
```

Or:

```
package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/aaronland/go-http-server"
	_ "github.com/aaronland/go-http-server-tsnet"
	"github.com/aaronland/go-http-server-tsnet/http/www"
)

func main() {

	server_uri := flag.String("server-uri", "tsnet://testing:80", "A valid aaronland/go-http-server URI.")
	flag.Parse()

	ctx := context.Background()
	s, _ := server.NewServer(ctx, *server_uri)

	handler := www.ExampleHandler()

	mux := http.NewServeMux()
	mux.Handle("/", handler)

	log.Printf("Listening on %s", s.Address())
	s.ListenAndServe(ctx, mux)
}
```

And:

```
package www

import (
	"fmt"
	"net/http"

	"github.com/aaronland/go-http-server-tsnet"
)

func ExampleHandler() http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		// tsnet.SetWhoIs is called/assigned by the middleware handler
		// implemented in tsnet.go
		
		who, _ := tsnet.GetWhoIs(req)

		login_name := who.UserProfile.LoginName
		computed_name := who.Node.ComputedName

		msg := fmt.Sprintf("Hello, %s (%s)", login_name, computed_name)
		rsp.Write([]byte(msg))
	}

	h := http.HandlerFunc(fn)
	return h
}
```

## See also

* https://github.com/aaronland/go-http-server
* https://pkg.go.dev/tailscale.com/tsnet
* https://tailscale.com/blog/tsnet-virtual-private-services/