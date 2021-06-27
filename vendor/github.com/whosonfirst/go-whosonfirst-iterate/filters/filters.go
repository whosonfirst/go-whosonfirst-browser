package filters

import (
	"context"
	"io"
)

type Filters interface {
	Apply(context.Context, io.ReadSeekCloser) (bool, error)
}
