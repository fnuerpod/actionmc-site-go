package authdatabase

import (
	dbh "git.dsrt-int.net/actionmc/actionmc-site-go/sqlite3dbh"
	_ "github.com/mattn/go-sqlite3"
)

func (db *MCAuthDB_sqlite3) Checkadmin(id string) (is_admin bool, privilege_level int) {
	db.logger.Debug.Println("Checking if user is an administrator (by their Discord ID)...")

	rows, err := db.handler.FetchRowsFromTable("admins", []dbh.FilterItem{
		{"discordID", dbh.CMP_EQ, id},
	})
	defer rows.Close()

	var uid string

	if err != nil {
		goto errorState
	}

	if rows.Next() {
		err = rows.Scan(&uid, &privilege_level)
	} else {
		is_admin = false
		privilege_level = 0
		return
	}

errorState:
	if err != nil {
		// shit went fuck, probably an invalid username or something.
		db.logger.Debug.Println("Failed to get user from database - not an administrator.")
		privilege_level = 0
		is_admin = false
	} else {
		db.logger.Debug.Println("Managed to get user from database - is an administrator.")
		is_admin = true
	}

	return

}
