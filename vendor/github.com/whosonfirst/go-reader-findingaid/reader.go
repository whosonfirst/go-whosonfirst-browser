package findingaid

import (
	_ "github.com/whosonfirst/go-reader-http"
)

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jtacoma/uritemplates"
	wof_reader "github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-reader-findingaid/finder"
	wof_uri "github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	"net/url"
)

// WHOSONFIRST_DATA_TEMPLATE is a URL template for the root `data` directory in Who's On First data repositories.
const WHOSONFIRST_DATA_TEMPLATE string = "https://raw.githubusercontent.com/whosonfirst-data/{repo}/master/data/"

// type FindingAidReader implements the `whosonfirst/go-reader` interface for use with Who's On First finding aids.
type FindingAidReader struct {
	wof_reader.Reader
	// A SQLite `sql.DB` instance containing Who's On First finding aid data.
	db *sql.DB
	// A compiled `uritemplates.UriTemplate` to use resolving Who's On First finding aid URIs.
	template *uritemplates.UriTemplate
	finder   finder.Finder
}

func init() {
	ctx := context.Background()
	wof_reader.RegisterReader(ctx, "findingaid", NewFindingAidReader)
}

// NewFindingAidReader will return a new `whosonfirst/go-reader.Reader` instance for reading Who's On First
// documents by first resolving a URL using a Who's On First finding aid.
func NewFindingAidReader(ctx context.Context, uri string) (wof_reader.Reader, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	fu := url.URL{}
	fu.Scheme = u.Host
	fu.Path = u.Path
	fu.RawQuery = u.RawQuery

	f_uri := fu.String()

	f, err := finder.NewFinder(ctx, f_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create finder, %w", err)
	}

	q := u.Query()

	uri_template := WHOSONFIRST_DATA_TEMPLATE

	if q.Get("template") != "" {
		uri_template = q.Get("template")
	}

	t, err := uritemplates.Parse(uri_template)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI template, %w", err)
	}

	r := &FindingAidReader{
		finder:   f,
		template: t,
	}

	return r, nil
}

// Read returns an `io.ReadSeekCloser` instance for the document resolved by `uri`.
func (r *FindingAidReader) Read(ctx context.Context, uri string) (io.ReadSeekCloser, error) {

	new_r, rel_path, err := r.getReaderAndPath(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive reader and path, %w", err)
	}

	return new_r.Read(ctx, rel_path)
}

// ReaderURI returns final URI resolved by `uri` for this reader.
func (r *FindingAidReader) ReaderURI(ctx context.Context, uri string) string {
	return uri
}

// getReaderAndPath returns a new `whosonfirst/go-reader.Reader` instance, and the relative path,
// to use for reading documents resolved by `uri`.
func (r *FindingAidReader) getReaderAndPath(ctx context.Context, uri string) (wof_reader.Reader, string, error) {

	reader_uri, rel_path, err := r.getReaderURIAndPath(ctx, uri)

	if err != nil {
		return nil, "", fmt.Errorf("Failed to derive reader URI and path, %w", err)
	}

	new_r, err := wof_reader.NewReader(ctx, reader_uri)

	if err != nil {
		return nil, "", fmt.Errorf("Failed to create reader, %w", err)
	}

	return new_r, rel_path, nil
}

// getReaderAndPath returns a new `whosonfirst/go-reader.Reader` URI, and the relative path,
// to use for reading documents resolved by `uri`.
func (r *FindingAidReader) getReaderURIAndPath(ctx context.Context, uri string) (string, string, error) {

	// TBD: cache this?

	id, uri_args, err := wof_uri.ParseURI(uri)

	if err != nil {
		return "", "", fmt.Errorf("Failed to parse URI, %w", err)
	}

	repo, err := r.finder.GetRepo(ctx, id)

	if err != nil {
		return "", "", fmt.Errorf("Failed to derive repo, %w", err)
	}

	rel_path, err := wof_uri.Id2RelPath(id, uri_args)

	if err != nil {
		return "", "", fmt.Errorf("Failed to derive path, %w", err)
	}

	values := map[string]interface{}{
		"repo": repo,
	}

	reader_uri, err := r.template.Expand(values)

	if err != nil {
		return "", "", fmt.Errorf("Failed to derive reader URI, %w", err)
	}

	return reader_uri, rel_path, nil
}
