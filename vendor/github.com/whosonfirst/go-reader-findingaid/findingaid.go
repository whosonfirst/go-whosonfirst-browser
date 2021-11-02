package findingaid

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jtacoma/uritemplates"
	_ "github.com/mattn/go-sqlite3"
	wof_reader "github.com/whosonfirst/go-reader"
	_ "github.com/whosonfirst/go-reader-http"
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

	q := u.Query()

	dsn := q.Get("dsn")

	db, err := sql.Open("sqlite3", dsn)

	if err != nil {
		return nil, fmt.Errorf("Failed to open database, %w", err)
	}

	uri_template := WHOSONFIRST_DATA_TEMPLATE

	if q.Get("template") != "" {
		uri_template = q.Get("template")
	}

	t, err := uritemplates.Parse(uri_template)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI template, %w", err)
	}

	f := &FindingAidReader{
		db:       db,
		template: t,
	}

	return f, nil
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

	repo, err := r.getRepo(ctx, id)

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

// getRepo returns the name of the repository associated with this ID in a Who's On First finding
// aid.
func (r *FindingAidReader) getRepo(ctx context.Context, id int64) (string, error) {

	var repo string

	q := "SELECT s.name FROM catalog c, sources s WHERE c.id = ? AND c.repo_id = s.id"

	row := r.db.QueryRowContext(ctx, q, id)
	err := row.Scan(&repo)

	if err != nil {
		return "", fmt.Errorf("Failed to scan row, %w", err)
	}

	return repo, nil
}
