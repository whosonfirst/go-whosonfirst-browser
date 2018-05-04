package tables

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/twpayne/go-geom"
	gogeom_geojson "github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkt"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/geometry"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-mysql"
	"github.com/whosonfirst/go-whosonfirst-mysql/utils"
	_ "log"
)

type WhosonfirstTable struct {
	mysql.Table
	name string
}

func NewWhosonfirstTableWithDatabase(db mysql.Database) (mysql.Table, error) {

	t, err := NewWhosonfirstTable()

	if err != nil {
		return nil, err
	}

	err = t.InitializeTable(db)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func NewWhosonfirstTable() (mysql.Table, error) {

	t := WhosonfirstTable{
		name: "whosonfirst",
	}

	return &t, nil
}

func (t *WhosonfirstTable) Name() string {
	return t.name
}

// https://dev.mysql.com/doc/refman/8.0/en/json-functions.html
// https://www.percona.com/blog/2016/03/07/json-document-fast-lookup-with-mysql-5-7/
// https://archive.fosdem.org/2016/schedule/event/mysql57_json/attachments/slides/1291/export/events/attachments/mysql57_json/slides/1291/MySQL_57_JSON.pdf

func (t *WhosonfirstTable) Schema() string {

	sql := `CREATE TABLE IF NOT EXISTS %s (
		      id BIGINT UNSIGNED PRIMARY KEY,
		      properties JSON NOT NULL,
		      geometry GEOMETRY NOT NULL,
		      centroid POINT NOT NULL COMMENT 'This is not necessary a math centroid',
		      lastmodified INT NOT NULL,
		      parent_id BIGINT       GENERATED ALWAYS AS (JSON_UNQUOTE(JSON_EXTRACT(properties,'$."wof:parent_id"'))) VIRTUAL,
		      placetype VARCHAR(64)  GENERATED ALWAYS AS (JSON_UNQUOTE(JSON_EXTRACT(properties,'$."wof:placetype"'))) VIRTUAL,
		      is_current TINYINT     GENERATED ALWAYS AS (JSON_UNQUOTE(JSON_EXTRACT(properties,'$."mz:is_current"'))) VIRTUAL,
		      is_nullisland TINYINT  GENERATED ALWAYS AS (JSON_LENGTH(JSON_EXTRACT(properties, '$."mz:is_nullisland"'))) VIRTUAL,
		      is_approximate TINYINT GENERATED ALWAYS AS (JSON_LENGTH(JSON_EXTRACT(properties, '$."mz:is_approximate"'))) VIRTUAL,
		      is_ceased TINYINT      GENERATED ALWAYS AS (JSON_UNQUOTE(JSON_EXTRACT(properties,'$."edtf:cessation"')) != "" AND json_unquote(json_extract(properties,'$."edtf:cessation"')) != "uuuu") VIRTUAL,
		      is_deprecated TINYINT  GENERATED ALWAYS AS (JSON_UNQUOTE(JSON_EXTRACT(properties,'$."edtf:deprecated"')) != "" AND json_unquote(json_extract(properties,'$."edtf:deprecated"')) != "uuuu") VIRTUAL,
		      is_superseded TINYINT  GENERATED ALWAYS AS (JSON_LENGTH(JSON_EXTRACT(properties, '$."wof:superseded_by"')) > 0) VIRTUAL,
		      is_superseding TINYINT GENERATED ALWAYS AS (JSON_LENGTH(JSON_EXTRACT(properties, '$."wof:supersedes"')) > 0) VIRTUAL,
		      date_upper DATE	     GENERATED ALWAYS AS (JSON_UNQUOTE(JSON_EXTRACT(properties, '$."date:cessation_upper"'))) VIRTUAL,
		      date_lower DATE	     GENERATED ALWAYS AS (JSON_UNQUOTE(JSON_EXTRACT(properties, '$."date:inception_lower"'))) VIRTUAL,
		      KEY parent_id (parent_id),
		      KEY placetype (placetype),
		      KEY is_current (is_current),
		      KEY is_nullisland (is_nullisland),
		      KEY is_approximate (is_approximate),
		      KEY is_deprecated (is_deprecated),
		      KEY is_superseded (is_superseded),
		      KEY is_superseding (is_superseding),
		      KEY date_upper (date_upper),
		      KEY date_lower (date_lower),
		      SPATIAL KEY idx_geometry (geometry),
		      SPATIAL KEY idx_centroid (centroid)
	      ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`

	// so while this is technically possible we're not doing it because it will return the math centroid
	// and well we all know what that means for places like SF - if you don't it means the math centroid
	// will be in the pacific ocean because technically the farralon islands are part of SF - so instead
	// we we compute the centroid below (using the go-whosonfirst-geojson-v2 Centroid interface)
	//
	// ALTER TABLE whosonfirst ADD centroid POINT GENERATED ALWAYS AS (ST_Centroid(geometry)) VIRTUAL;
	// (20180504/thisisaaronland)

	return fmt.Sprintf(sql, t.Name())
}

func (t *WhosonfirstTable) InitializeTable(db mysql.Database) error {

	return utils.CreateTableIfNecessary(db, t)
}

func (t *WhosonfirstTable) IndexRecord(db mysql.Database, i interface{}) error {
	return t.IndexFeature(db, i.(geojson.Feature))
}

func (t *WhosonfirstTable) IndexFeature(db mysql.Database, f geojson.Feature) error {

	conn, err := db.Conn()

	if err != nil {
		return err
	}

	tx, err := conn.Begin()

	if err != nil {
		return err
	}

	str_geom, err := geometry.ToString(f)

	if err != nil {
		return err
	}

	var g geom.T
	err = gogeom_geojson.Unmarshal([]byte(str_geom), &g)

	if err != nil {
		return err
	}

	wkt_geom, err := wkt.Marshal(g)

	centroid, err := whosonfirst.Centroid(f)

	if err != nil {
		return err
	}

	coord := centroid.Coord()

	wkt_centroid := fmt.Sprintf("POINT(%0.6f %0.6f)", coord.X, coord.Y)

	sql := fmt.Sprintf(`REPLACE INTO %s (
		geometry, centroid, id, properties, lastmodified
	) VALUES (
		ST_GeomFromText('%s'), ST_GeomFromText('%s'), ?, ?, ?
	)`, t.Name(), wkt_geom, wkt_centroid)

	stmt, err := tx.Prepare(sql)

	if err != nil {
		return err
	}

	defer stmt.Close()

	props := gjson.GetBytes(f.Bytes(), "properties")
	props_json, err := json.Marshal(props.Value())

	if err != nil {
		return err
	}

	lastmod := whosonfirst.LastModified(f)

	_, err = stmt.Exec(f.Id(), string(props_json), lastmod)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
