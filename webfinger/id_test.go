package webfinger

import (
	"testing"
)

func TestDeriveWhosOnFirstURI(t *testing.T) {

	tests_ok := map[string]string{
		"acct:85922583@whosonfirst.org": "85922583",
	}

	tests_fail := []string{
		"85922583",
		"85922583@whosonfirst.org",
	}

	for resource, expected := range tests_ok {

		wof_uri, err := DeriveWhosOnFirstURIFromResource(resource)

		if err != nil {
			t.Fatalf("Failed to parse resource '%s', %v", resource, err)
		}

		if wof_uri != expected {
			t.Fatalf("Unexpected result for '%s'. Expected '%s' but got '%s'.", resource, expected, wof_uri)
		}
	}

	for _, resource := range tests_fail {

		_, err := DeriveWhosOnFirstURIFromResource(resource)

		if err == nil {
			t.Fatalf("Expected '%s' to fail but did not.", resource)
		}
	}

}
