package custom

import (
	"context"
)

type CustomValidationFunc func(context.Context, []byte) error
