package webfinger

import (
	"fmt"
	"net/mail"
	"strings"
)

const AccountScheme string = "acct:"

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
