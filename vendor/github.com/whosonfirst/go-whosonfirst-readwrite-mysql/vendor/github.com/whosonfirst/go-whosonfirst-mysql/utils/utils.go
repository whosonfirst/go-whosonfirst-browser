package utils

import (
	_ "database/sql"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-mysql"
)

var lookup_table map[string]bool

func init() {
	lookup_table = make(map[string]bool)
}

func HasTable(db mysql.Database, table string) (bool, error) {

	dsn := db.DSN()

	lookup_key := fmt.Sprintf("%s#%s", dsn, table)

	has_table, ok := lookup_table[lookup_key]

	if ok {
		return has_table, nil
	}

	conn, err := db.Conn()

	if err != nil {
		return false, err
	}

	// Would that the following work in Go... because it totally works
	// from the MySQL CLI... I have no idea (20180426/thisisaaronland)
	// 2018/04/26 09:37:45 ERR SHOW TABLES LIKE ?
	// Error 1064: You have an error in your SQL syntax; check the manual that corresponds to your MariaDB server version for the right syntax to use near '?' at line 1

	/*
		query := "SHOW TABLES LIKE ?"
		row := conn.QueryRow(query, table)
		err = row.Scan()

		log.Println("ERR", query, err)

		switch {
		case err == sql.ErrNoRows:
			return false, nil
		case err != nil:
			return false, err
		default:
			return true, nil
		}
	*/

	has_table = false

	query := "SHOW TABLES"
	rows, err := conn.Query(query)

	if err != nil {
		return false, err
	}

	defer rows.Close()

	for rows.Next() {

		var name string
		err := rows.Scan(&name)

		if err != nil {
			return false, err
		}

		if name == table {
			has_table = true
			break
		}
	}

	lookup_table[lookup_key] = has_table

	return has_table, nil
}

func CreateTableIfNecessary(db mysql.Database, t mysql.Table) error {

	create := false

	if db.DSN() == ":memory:" {
		create = true
	} else {

		has_table, err := HasTable(db, t.Name())

		if err != nil {
			return err
		}

		if !has_table {
			create = true
		}
	}

	if create {

		sql := t.Schema()

		conn, err := db.Conn()

		if err != nil {
			return err
		}

		_, err = conn.Exec(sql)

		if err != nil {
			return err
		}

	}

	return nil
}
