package html

import (
	"context"
	"embed"
	sfom_html "github.com/sfomuseum/go-template/html"
	"html/template"
)

//go:embed *.html
var FS embed.FS

func LoadTemplates(ctx context.Context) (*template.Template, error) {

	return sfom_html.LoadTemplates(ctx, FS)
}
