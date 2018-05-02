package reader

import (
	"errors"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-readwrite/bytes"
	wof_reader "github.com/whosonfirst/go-whosonfirst-readwrite/reader"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	"github.com/whosonfirst/go-whosonfirst-sqlite-features/tables"
	"github.com/whosonfirst/go-whosonfirst-sqlite/database"
	"github.com/whosonfirst/go-whosonfirst-sqlite/utils"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
)

type SQLiteReader struct {
	wof_reader.Reader
	database *database.SQLiteDatabase
	table    sqlite.Table
}

func NewSQLiteReader(dsn string, args ...interface{}) (wof_reader.Reader, error) {

	db, err := database.NewDBWithDriver("sqlite3", dsn)

	if err != nil {
		return nil, err
	}

	tbl, err := tables.NewGeoJSONTable()

	if err != nil {
		return nil, err
	}

	ok, err := utils.HasTable(db, tbl.Name())

	if err != nil {
		return nil, err
	}

	if ok == false {
		return nil, errors.New(fmt.Sprintf("Database is missing %s table", tbl.Name()))
	}

	r := SQLiteReader{
		database: db,
		table:    tbl,
	}

	return &r, nil
}

func (r *SQLiteReader) Read(path string) (io.ReadCloser, error) {

	id, err := uri.IdFromPath(path)

	if err != nil {
		return nil, err
	}

	conn, err := r.database.Conn()

	if err != nil {
		return nil, err
	}

	q := fmt.Sprintf("SELECT body FROM %s WHERE id=?", r.table.Name())
	row := conn.QueryRow(q, id)

	var body string
	err = row.Scan(&body)

	if err != nil {
		return nil, err
	}

	return bytes.ReadCloserFromBytes([]byte(body))
}

func (r *SQLiteReader) URI(path string) string {
     return fmt.Sprintf("sqlite://%s/%s#%s", r.database.DSN(), r.table.Name(), path)
}
