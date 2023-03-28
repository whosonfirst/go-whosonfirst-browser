package chrome

import (
	"context"
	_ "log"
	"net/http"
)

func init() {
	ctx := context.Background()
	RegisterChrome(ctx, "none", NewNoneChrome)
}

// type NoneChrome implements the Chrome interface that always returns a "not authorized" error.
type NoneChrome struct {
	Chrome
}

// NewNoneChrome implements the Chrome interface that always returns a "not authorized" error.
// configured by 'uri' which is expected to take the form of:
//
//	none://
func NewNoneChrome(ctx context.Context, uri string) (Chrome, error) {
	c := &NoneChrome{}
	return c, nil
}

// WrapHandler returns 'h' unchanged.
func (c *NoneChrome) WrapHandler(h http.Handler, path string) http.Handler {
	return h
}

func (c *NoneChrome) AppendStaticAssetHandlers(*http.ServeMux) error {
	return nil
}

func (c *NoneChrome) AppendStaticAssetHandlersWithPrefix(mux *http.ServeMux, prefix string) error {
	return nil
}
