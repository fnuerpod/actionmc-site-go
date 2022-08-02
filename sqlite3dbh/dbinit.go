// SQLITE3DBH is a mini library built for AMC-site-go.
// This librarys primary function is to provide a safer and
// more programatic way to work with sqlite databases.
package sqlite3dbh

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// OpenFromFile opens a sqlite3 database from a file
// An error is passed from sql.Open if there is any
func OpenFromFile(name string) (*DBHandler, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	} else {
		return &DBHandler{db}, nil
	}

}

// OpenFromDB creates a DBHandler from an already init sql.DB
// It assumes the sql.DB uses the github.com/mattn/go-sqlite3 implementation
func OpenFromDB(db *sql.DB) *DBHandler {
	return &DBHandler{db}
}
