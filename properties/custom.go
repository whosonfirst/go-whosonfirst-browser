package properties

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
	"github.com/tidwall/gjson"
)

type CustomProperty interface {
	Name() string
	Type() string
	Required() bool
}

func EnsureCustomPropertyHasValue(ctx context.Context, pr CustomProperty, body []byte) (bool, error) {

	if !pr.Required() {
		return true, nil
	}

	path := fmt.Sprintf("properties.%s", pr.Name())

	rsp := gjson.GetBytes(body, path)

	if !rsp.Exists() {
		return false, nil
	}

	if rsp.String() == "" {
		return false, nil
	}

	return true, nil
}

var custom_roster roster.Roster

// CustomInitializationFunc is a function defined by individual custom package and used to create
// an instance of that custom
type CustomPropertyInitializationFunc func(ctx context.Context, uri string) (CustomProperty, error)

// RegisterCustom registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `CustomProperty` instances by the `NewCustom` method.
func RegisterCustomProperty(ctx context.Context, scheme string, init_func CustomPropertyInitializationFunc) error {

	err := ensureCustomRoster()

	if err != nil {
		return err
	}

	return custom_roster.Register(ctx, scheme, init_func)
}

func ensureCustomRoster() error {

	if custom_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		custom_roster = r
	}

	return nil
}

// NewCustom returns a new `CustomProperty` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `CustomInitializationFunc`
// function used to instantiate the new `Custom`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterCustom` method.
func NewCustomProperty(ctx context.Context, uri string) (CustomProperty, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := custom_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(CustomPropertyInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureCustomRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range custom_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}