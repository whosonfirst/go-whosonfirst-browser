package properties

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

const CUSTOM_STRING_PROPERTY string = "string"

type CustomStringProperty struct {
	CustomProperty
	name           string
	required       bool
	custom_element string
}

func init() {
	ctx := context.Background()
	RegisterCustomProperty(ctx, CUSTOM_STRING_PROPERTY, NewCustomStringProperty)
}

func NewCustomStringProperty(ctx context.Context, uri string) (CustomProperty, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	name := q.Get("name")

	if name == "" {
		return nil, fmt.Errorf("Missing ?name= parameter")
	}

	pr := &CustomStringProperty{
		name: name,
	}

	str_required := q.Get("required")

	if str_required != "" {

		r, err := strconv.ParseBool(str_required)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?required= parameter, %w", err)
		}

		pr.required = r
	}

	custom_el := q.Get("custom-element")

	if custom_el != "" {

		// Validate custom_el here...
		pr.custom_element = custom_el
	}

	return pr, nil
}

func (pr *CustomStringProperty) Name() string {
	return pr.name
}

func (pr *CustomStringProperty) Required() bool {
	return pr.required
}

func (pr *CustomStringProperty) Type() string {
	return CUSTOM_STRING_PROPERTY
}

func (pr *CustomStringProperty) CustomElement() string {
	return pr.custom_element
}
