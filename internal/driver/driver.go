// This file can hold the implementations when connecting to other dbs besides postgres
package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifeTime = 5 * time.Minute

// Connects sql and creates creates db pool for PG
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbLifeTime)

	dbConn.SQL = d

	err = TestDb(dbConn.SQL)
	if err != nil {
		return dbConn, nil
	}

	return dbConn, nil
}

// Creates new db instance
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn) // Tells our app we are using pgx database driver
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// Test db connection, pings db
func TestDb(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}

	return nil
}
