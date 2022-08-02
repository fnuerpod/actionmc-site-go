package authdatabase

import (
	//dbh "git.dsrt-int.net/actionmc/actionmc-site-go/sqlite3dbh"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func (db *MCAuthDB_sqlite3) GetPhotos() (rows *sql.Rows) {
	db.logger.Debug.Println("Getting all photos from database...")
	rows, err := db.handler.FetchTable("site_gallery")
	if err != nil {
		log.Fatal(err)
	}
	db.logger.Debug.Println("Got all photos from database.")
	return
}
