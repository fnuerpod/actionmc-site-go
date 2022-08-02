package authdatabase

import (
	"errors"
	dbh "git.dsrt-int.net/actionmc/actionmc-site-go/sqlite3dbh"
	_ "github.com/mattn/go-sqlite3"
)

func (db *MCAuthDB_sqlite3) Getstate(id string) (state string, ok bool) {
	db.logger.Debug.Println("Getting user state from database (by their Discord ID)...")

	var uid, una, ena string

	stmt, err := db.handler.FetchRowsFromTableStmt("users", []dbh.FilterItem{
		{"discordID", dbh.CMP_EQ, nil},
	})

	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&uid, &una, &state, &ena)

	if err != nil {
		// shit went fuck, probably an invalid username or something.
		ok = false
		db.logger.Debug.Println("Failed to obtain user state from database. Invalid user(?)")
	} else {
		ok = true
		db.logger.Debug.Println("Managed to obtain user from database.")
	}

	return

}

func (db *MCAuthDB_sqlite3) Getstate_byuname(uname string) (state string, ok bool) {
	// start database transaction
	db.logger.Debug.Println("Getting user state from database (by their username on Minecraft)...")

	var uid, una, ena string

	rows, err := db.handler.FetchRowsFromTable("users", []dbh.FilterItem{
		{"userName", dbh.CMP_EQ, uname},
	})
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&uid, &una, &state, &ena)
	} else {
		err = errors.New("No user")
	}

	if err != nil {
		// shit went fuck, probably an invalid username or something.
		ok = false
		db.logger.Debug.Println("Failed to obtain user state from database. Invalid user(?)")
	} else {
		db.logger.Debug.Println("Managed to obtain user from database.")
		ok = true
	}

	return

}
