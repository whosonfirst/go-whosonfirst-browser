package reader

import (
	"errors"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-mysql"
	"github.com/whosonfirst/go-whosonfirst-mysql/database"
	"github.com/whosonfirst/go-whosonfirst-mysql/tables"
	"github.com/whosonfirst/go-whosonfirst-mysql/utils"
	reader_utils "github.com/whosonfirst/go-whosonfirst-readwrite-mysql/utils"
	"github.com/whosonfirst/go-whosonfirst-readwrite/bytes"
	wof_reader "github.com/whosonfirst/go-whosonfirst-readwrite/reader"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	"log"
)

type MySQLGeoJSONReader struct {
	wof_reader.Reader
	database *database.MySQLDatabase
	table    mysql.Table
}

func NewMySQLGeoJSONReader(dsn string, args ...interface{}) (wof_reader.Reader, error) {

	db, err := database.NewDB(dsn)

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

	r := MySQLGeoJSONReader{
		database: db,
		table:    tbl,
	}

	return &r, nil
}

func (r *MySQLGeoJSONReader) Read(path string) (io.ReadCloser, error) {

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

	var str_body string
	err = row.Scan(&str_body)

	if err != nil {
		return nil, err
	}

	log.Println(str_body)

	return bytes.ReadCloserFromBytes([]byte(str_body))
}

func (r *MySQLGeoJSONReader) URI(path string) string {
	return fmt.Sprintf("mysql://%s/%s#%s", reader_utils.ScrubDSN(r.database.DSN()), r.table.Name(), path)
}
