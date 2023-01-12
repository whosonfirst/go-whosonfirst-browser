package webfinger

import (
	"fmt"
	"net/mail"
	"strings"
)

// AccountScheme is URI scheme for WebFinger `acct://` URIs.
const AccountScheme string = "acct:"

// DeriveWhosOnFirstURIFromResource derives a valid Who's On First URI string (something that can be parsed by
// the `whosonfirst/go-whosonfirst-uri` package from 'resource'.
func DeriveWhosOnFirstURIFromResource(resource string) (string, error) {

	if !strings.HasPrefix(resource, AccountScheme) {
		return "", fmt.Errorf("URI is missing account scheme")
	}

	str_addr := strings.Replace(resource, AccountScheme, "", 1)

	addr, err := mail.ParseAddress(str_addr)

	if err != nil {
		return "", fmt.Errorf("Failed to parse address, %w", err)
	}

	parts := strings.Split(addr.Address, "@")

	return parts[0], nil
}
