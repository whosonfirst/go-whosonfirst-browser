package reader

import (
	"encoding/json"
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
)

type MySQLWhosonfirstReader struct {
	wof_reader.Reader
	database *database.MySQLDatabase
	table    mysql.Table
}

type Feature struct {
	Type       string      `json:"type"`
	Geometry   interface{} `json:"geometry"`
	Properties interface{} `json:"properties"`
}

func NewMySQLWhosonfirstReader(dsn string, args ...interface{}) (wof_reader.Reader, error) {

	db, err := database.NewDB(dsn)

	if err != nil {
		return nil, err
	}

	tbl, err := tables.NewWhosonfirstTable()

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

	r := MySQLWhosonfirstReader{
		database: db,
		table:    tbl,
	}

	return &r, nil
}

func (r *MySQLWhosonfirstReader) Read(path string) (io.ReadCloser, error) {

	id, err := uri.IdFromPath(path)

	if err != nil {
		return nil, err
	}

	conn, err := r.database.Conn()

	if err != nil {
		return nil, err
	}

	q := fmt.Sprintf("SELECT ST_AsGeoJSON(geometry), JSON_UNQUOTE(properties) FROM %s WHERE id=?", r.table.Name())
	row := conn.QueryRow(q, id)

	var str_geom string
	var str_props string

	err = row.Scan(&str_geom, &str_props)

	if err != nil {
		return nil, err
	}

	var geom interface{}
	var props interface{}

	err = json.Unmarshal([]byte(str_geom), &geom)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(str_props), &props)

	if err != nil {
		return nil, err
	}

	f := Feature{
		Type:       "feature",
		Geometry:   geom,
		Properties: props,
	}

	body, err := json.Marshal(f)

	if err != nil {
		return nil, err
	}

	return bytes.ReadCloserFromBytes(body)
}

func (r *MySQLWhosonfirstReader) URI(path string) string {
	return fmt.Sprintf("mysql://%s/%s#%s", reader_utils.ScrubDSN(r.database.DSN()), r.table.Name(), path)
}
