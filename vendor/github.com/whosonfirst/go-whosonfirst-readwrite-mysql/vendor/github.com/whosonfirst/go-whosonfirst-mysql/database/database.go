package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "log"
	"sync"
)

type MySQLDatabase struct {
	conn *sql.DB
	dsn  string
	mu   *sync.Mutex
}

func NewDB(dsn string) (*MySQLDatabase, error) {

	conn, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	mu := new(sync.Mutex)

	db := MySQLDatabase{
		conn: conn,
		dsn:  dsn,
		mu:   mu,
	}

	return &db, err
}

func (db *MySQLDatabase) Lock() {
	db.mu.Lock()
}

func (db *MySQLDatabase) Unlock() {
	db.mu.Unlock()
}

func (db *MySQLDatabase) Conn() (*sql.DB, error) {
	return db.conn, nil
}

func (db *MySQLDatabase) DSN() string {
	return db.dsn
}

func (db *MySQLDatabase) Close() error {
	return db.conn.Close()
}
