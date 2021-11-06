package resolver

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/url"
)

// type SQLiteResolver implements the `Resolver` interface for data stored in a SQLite database..
type SQLiteResolver struct {
	Resolver
	// A SQLite `sql.DB` instance containing Who's On First finding aid data.
	db *sql.DB
}

func init() {
	ctx := context.Background()
	RegisterResolver(ctx, "sqlite", NewSQLiteResolver)
	RegisterResolver(ctx, "sqlite3", NewSQLiteResolver)
}

// NewSQLiteResolver will return a new `Resolver` instance for resolving repository names
// and IDs stored in a SQLite database.
func NewSQLiteResolver(ctx context.Context, uri string) (Resolver, error) {

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

	f := &SQLiteResolver{
		db: db,
	}

	return f, nil
}

// GetRepo returns the name of the repository associated with this ID in a Who's On First finding aid.
func (r *SQLiteResolver) GetRepo(ctx context.Context, id int64) (string, error) {

	var repo string

	q := "SELECT s.name FROM catalog c, sources s WHERE c.id = ? AND c.repo_id = s.id"

	row := r.db.QueryRowContext(ctx, q, id)
	err := row.Scan(&repo)

	if err != nil {
		return "", fmt.Errorf("Failed to scan row, %w", err)
	}

	return repo, nil
}
