package writer

import (
	"errors"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-mysql"
	"github.com/whosonfirst/go-whosonfirst-mysql/database"
	"github.com/whosonfirst/go-whosonfirst-mysql/tables"
	wof_writer "github.com/whosonfirst/go-whosonfirst-readwrite/writer"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
)

type SQLiteWriter struct {
	wof_writer.Writer
	database *database.MySQLDatabase
	table    mysql.Table
}

func NewSQLiteWriter(dsn string, args ...interface{}) (wof_writer.Writer, error) {

	db, err := database.NewDB(dsn)

	if err != nil {
		return nil, err
	}

	tbl, err := tables.NewWhosonfirstTableWithDatabase(db)

	if err != nil {
		return nil, err
	}

	wr := SQLiteWriter{
		database: db,
		table:    tbl,
	}

	return &wr, nil
}

func (wr *SQLiteWriter) Write(path string, fh io.ReadCloser) error {

	id, err := uri.IdFromPath(path)

	if err != nil {
		return err
	}

	return errors.New(fmt.Sprintf("Please write %d (%s) to the database", id, path))
}

func (wr *SQLiteWriter) URI(path string) string {
	return fmt.Sprintf("sqlite://%s/%s#%s", wr.database.DSN(), wr.table.Name(), path)
}
