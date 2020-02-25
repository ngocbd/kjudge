package db

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

// DBContext defines a slim common interface between Db and Tx.
type DBContext interface {
	// Get binds a single row into a struct.
	Get(result interface{}, query string, args ...interface{}) error
	// Select binds rows into an array of structs.
	Select(results interface{}, query string, args ...interface{}) error
	// Exec executes a query.
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// DB is an implementation of the underlying DB.
type DB struct {
	*sqlx.DB

	PersistentConn *sqlite3.SQLiteConn
}

// New creates a new DB object from the given filename.
func New(filename string) (*DB, error) {
	dsn := fmt.Sprintf("%s?_fk=1&mode=rw", filename)
	sqlxdb, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	conn, err := sqlxdb.Driver().Open(dsn)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	db := &DB{
		DB:             sqlxdb,
		PersistentConn: conn.(*sqlite3.SQLiteConn),
	}
	// Perform migrations, if needed.
	if err := db.migrate(); err != nil {
		return nil, err
	}

	return db, nil
}
