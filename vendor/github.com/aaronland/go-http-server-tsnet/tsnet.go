package tsnet

// https://tailscale.com/blog/tsnet-virtual-private-services/
// https://github.com/tailscale/tailscale/blob/v1.32.2/tsnet/example/tshello/tshello.go
// https://pkg.go.dev/tailscale.com/tsnet
// https://tailscale.com/kb/1244/tsnet/

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/aaronland/go-http-server"
	"tailscale.com/client/tailscale"
	"tailscale.com/client/tailscale/apitype"
	"tailscale.com/tsnet"
)

// WHOIS_CONTEXT_KEY is the key used to store `apitype.WhoIsResponse` instances in a `context.Context` instance.
const WHOIS_CONTEXT_KEY string = "x-urn:aaronland:tsnet#whois"

func init() {
	ctx := context.Background()
	server.RegisterServer(ctx, "tsnet", NewTSNetServer)
}

// TSNServer implements the `Server` interface for a `net/http` server running as a Tailscale virtual private service.
type TSNetServer struct {
	server.Server
	tsnet_server *tsnet.Server
	hostname     string
	port         string
	funnel       bool
}

// NewTSNetServer returns a new `TSNetServer` instance configured by 'uri' which is
// expected to be defined in the form of:
//
//	tsnet://{HOSTNAME}:{PORT}?{PARAMETERS}
//
// Valid parameters are:
// * `auth-key` is a valid Tailscale auth key. If absent it is assumed that a valid `TS_AUTH_KEY` environment variable has already been set.
func NewTSNetServer(ctx context.Context, uri string) (server.Server, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	hostname := u.Hostname()
	port := u.Port()

	q := u.Query()

	tsnet_server := &tsnet.Server{
		Hostname: u.Hostname(),
	}

	auth_key := q.Get("auth-key")

	if auth_key != "" {

		// I don't really understand this...
		// 2022/11/06 10:17:17 Authkey is set; but state is NoState. Ignoring authkey. Re-run with TSNET_FORCE_LOGIN=1 to force use of authkey.

		err := os.Setenv("TSNET_FORCE_LOGIN", "1")

		if err != nil {
			return nil, fmt.Errorf("Failed to set TSNET_FORCE_LOGIN environment variable, %w", err)
		}

		tsnet_server.AuthKey = auth_key

	}

	q_ephemeral := q.Get("ephemeral")

	if q_ephemeral != "" {

		v, err := strconv.ParseBool(q_ephemeral)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?ephemeral query parameter, %w", err)
		}

		tsnet_server.Ephemeral = v
	}

	funnel := false

	q_funnel := q.Get("funnel")

	if q_funnel != "" {

		v, err := strconv.ParseBool(q_funnel)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?funnel query parameter, %w", err)
		}

		funnel = v
	}

	s := &TSNetServer{
		tsnet_server: tsnet_server,
		hostname:     hostname,
		port:         port,
		funnel:       funnel,
	}

	return s, nil
}

// ListenAndServe starts the server and listens for requests using 'mux' for routing. Additionally each handler in mux
// will be wrapped by a middleware handler that will ensure a Tailscale `api.WhoIsResponse` instance can be derived from
// the current request and then store that value in the request's context. This value can be retrieved using the `GetWhoIs`
// method.
func (s *TSNetServer) ListenAndServe(ctx context.Context, mux http.Handler) error {

	// It is important to include the hostname here
	addr := fmt.Sprintf(":%s", s.port)

	var listener net.Listener

	if s.funnel {

		l, err := s.tsnet_server.ListenFunnel("tcp", ":443")

		if err != nil {
			return fmt.Errorf("Failed to create funnel listener, %w", err)
		}

		listener = l
	} else {

		l, err := s.tsnet_server.Listen("tcp", addr)

		if err != nil {
			return fmt.Errorf("Failed to announce server, %w", err)
		}

		// https://testing
		// 2022/11/06 11:01:23 http: TLS handshake error from a.b.c.d:61940: 400 Bad Request: invalid domain

		// https://{TAILSCALE_IP}
		// 2022/11/06 11:01:46 http: TLS handshake error from a.b.c.d:53770: no SNI ServerName

		if s.port == "443" {

			l = tls.NewListener(l, &tls.Config{
				GetCertificate: tailscale.GetCertificate,
			})
		}

		listener = l
	}

	defer listener.Close()

	lc, err := s.tsnet_server.LocalClient()

	if err != nil {
		return fmt.Errorf("Failed to create local client, %w", err)
	}

	who_wrapper := func(next http.Handler) http.Handler {

		fn := func(rsp http.ResponseWriter, req *http.Request) {

			who, err := lc.WhoIs(req.Context(), req.RemoteAddr)

			if err != nil {
				http.Error(rsp, err.Error(), 500)
				return
			}

			who_req := SetWhoIs(req, who)
			next.ServeHTTP(rsp, who_req)
		}

		return http.HandlerFunc(fn)
	}

	err = http.Serve(listener, who_wrapper(mux))

	if err != nil {
		return fmt.Errorf("Failed to serve requests, %w", err)
	}

	return nil
}

// Address returns the fully-qualified URI where the server instance can be contacted.
func (s *TSNetServer) Address() string {

	var address string

	if s.port == "443" {
		address = fmt.Sprintf("https://%s", s.hostname)
	} else {
		address = fmt.Sprintf("http://%s:%s", s.hostname, s.port)
	}

	return address
}

// SetWhoIs will store 'who' in 'req.Context'.
func SetWhoIs(req *http.Request, who *apitype.WhoIsResponse) *http.Request {

	ctx := req.Context()

	who_ctx := context.WithValue(ctx, WHOIS_CONTEXT_KEY, who)
	who_req := req.WithContext(who_ctx)

	return who_req
}

// GetWhoIs will return the Tailscale `apitype.WhoIsResponse` instance stored in the 'req.Context'.
func GetWhoIs(req *http.Request) (*apitype.WhoIsResponse, error) {

	ctx := req.Context()

	v := ctx.Value(WHOIS_CONTEXT_KEY)

	if v == nil {
		return nil, fmt.Errorf("Unable to determine whois context")
	}

	switch v.(type) {
	case *apitype.WhoIsResponse:
		// pass
	default:
		return nil, fmt.Errorf("Invalid whois context")
	}

	return v.(*apitype.WhoIsResponse), nil
}
