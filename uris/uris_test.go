package uris

import (
	"net/url"
	"testing"
)

func TestURIs(t *testing.T) {

	u := &URIs{
		Route: "/id",
	}

	prefix := "/places"

	expected, err := url.JoinPath(prefix, u.Route)

	if err != nil {
		t.Fatal(err)
	}

	err = u.ApplyPrefix(prefix)

	if err != nil {
		t.Fatal(err)
	}

	if u.Route != expected {
		t.Fatalf("Expected '%s' but got '%s'", expected, u.Route)
	}
}
