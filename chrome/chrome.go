// package chrome provides a simple interface for enforcing chromeentication in HTTP handlers.
package chrome

import (
	"context"
	"fmt"
	_ "log"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

// type Chrome is a simple interface for	enforcing chromeentication in HTTP handlers.
type Chrome interface {
	// WrapHandler wraps a `http.Handler` with any implementation-specific middleware.
	WrapHandler(http.Handler, string) http.Handler
	AppendStaticAssetHandlers(*http.ServeMux) error
	AppendStaticAssetHandlersWithPrefix(*http.ServeMux, string) error
}

var chrome_roster roster.Roster

// ChromeInitializationFunc is a function defined by individual chrome package and used to create
// an instance of that chrome
type ChromeInitializationFunc func(ctx context.Context, uri string) (Chrome, error)

// RegisterChrome registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `Chrome` instances by the `NewChrome` method.
func RegisterChrome(ctx context.Context, scheme string, init_func ChromeInitializationFunc) error {

	err := ensureChromeRoster()

	if err != nil {
		return err
	}

	return chrome_roster.Register(ctx, scheme, init_func)
}

func ensureChromeRoster() error {

	if chrome_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		chrome_roster = r
	}

	return nil
}

// NewChrome returns a new `Chrome` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `ChromeInitializationFunc`
// function used to instantiate the new `Chrome`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterChrome` method.
func NewChrome(ctx context.Context, uri string) (Chrome, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := chrome_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(ChromeInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureChromeRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range chrome_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
