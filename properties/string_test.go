package properties

import (
	"context"
	"testing"
)

func TestCustomStringProperty(t *testing.T) {

	ctx := context.Background()

	pr, err := NewCustomProperty(ctx, "string://?name=sfo:level&required=true")

	if err != nil {
		t.Fatalf("Failed to create new custom string property, %v", err)
	}

	if pr.Name() != "sfo:level" {
		t.Fatalf("Unexpected name, %s", pr.Name())
	}

	if !pr.Required() {
		t.Fatalf("Expected custom property to be required")
	}

	if pr.Type() != CUSTOM_STRING_PROPERTY {
		t.Fatalf("Unexpected property type: %s", pr.Type())
	}
}
