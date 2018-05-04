package tables

import (
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-mysql"
	"github.com/whosonfirst/go-whosonfirst-mysql/utils"
	_ "log"
)

type GeoJSONTable struct {
	mysql.Table
	name string
}

func NewGeoJSONTableWithDatabase(db mysql.Database) (mysql.Table, error) {

	t, err := NewGeoJSONTable()

	if err != nil {
		return nil, err
	}

	err = t.InitializeTable(db)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func NewGeoJSONTable() (mysql.Table, error) {

	t := GeoJSONTable{
		name: "geojson",
	}

	return &t, nil
}

func (t *GeoJSONTable) Name() string {
	return t.name
}

func (t *GeoJSONTable) Schema() string {

	sql := `CREATE TABLE IF NOT EXISTS %s (
		      id BIGINT UNSIGNED PRIMARY KEY,
		      body LONGBLOB NOT NULL,
		      lastmodified INT NOT NULL,
		      KEY lastmodified (lastmodified)
	      ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`

	return fmt.Sprintf(sql, t.Name())
}

func (t *GeoJSONTable) InitializeTable(db mysql.Database) error {

	return utils.CreateTableIfNecessary(db, t)
}

func (t *GeoJSONTable) IndexRecord(db mysql.Database, i interface{}) error {
	return t.IndexFeature(db, i.(geojson.Feature))
}

func (t *GeoJSONTable) IndexFeature(db mysql.Database, f geojson.Feature) error {

	conn, err := db.Conn()

	if err != nil {
		return err
	}

	tx, err := conn.Begin()

	if err != nil {
		return err
	}

	sql := fmt.Sprintf(`REPLACE INTO %s (
		id, body, lastmodified
	) VALUES (
		?, ?, ?
	)`, t.Name())

	stmt, err := tx.Prepare(sql)

	if err != nil {
		return err
	}

	defer stmt.Close()

	if err != nil {
		return err
	}

	body := f.Bytes()
	lastmod := whosonfirst.LastModified(f)

	_, err = stmt.Exec(f.Id(), string(body), lastmod)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
