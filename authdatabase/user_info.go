package authdatabase

import (
	//"bufio"
	//"bytes"
	//"crypto/sha256"
	//"database/sql"
	//"errors"
	//"git.dsrt-int.net/actionmc/actionmc-site-go/logging"
	dbh "git.dsrt-int.net/actionmc/actionmc-site-go/sqlite3dbh"
	_ "github.com/mattn/go-sqlite3"
	//"log"
	//"strconv"
)

func (db *MCAuthDB_sqlite3) Getuser(id string) (UserInfo, bool) {
	db.logger.Debug.Println("Getting user from database (by their Discord ID)...")

	stmt, err := db.handler.FetchRowsFromTableStmt("users", []dbh.FilterItem{
		{"discordID", dbh.CMP_EQ, nil},
	})

	defer stmt.Close()

	var uid, name, state string
	var userChange int

	if err != nil {
		goto errorState
	}

	err = stmt.QueryRow(id).Scan(&uid, &name, &state, &userChange)

errorState:
	if err != nil {

		return UserInfo{}, false

		db.logger.Debug.Println("Failed to obtain user from database. Invalid user(?)")
	} else {
		db.logger.Debug.Println("Managed to obtain user from database.")
	}

	return UserInfo{
		Uid:        uid,
		Name:       name,
		State:      state,
		UserChange: userChange,
	}, true

}

func (db *MCAuthDB_sqlite3) GetTOTPSecret(id string) (string, bool) {
	db.logger.Debug.Println("Getting TOTP secret from database (by Discord ID)...")

	stmt, err := db.handler.FetchRowsFromTableStmt("sult_keys", []dbh.FilterItem{
		{"discordID", dbh.CMP_EQ, nil},
	})

	defer stmt.Close()

	var uid, key string

	if err != nil {
		goto errorState
	}

	err = stmt.QueryRow(id).Scan(&uid, &key)

errorState:
	if err != nil {

		return "", false

		db.logger.Debug.Println("Failed to obtain key from database. Invalid ID(?)")
	} else {
		db.logger.Debug.Println("Managed to obtain key from database.")
	}

	return key, true

}

func (db *MCAuthDB_sqlite3) Getuser_byuname(username string) (UserInfo, bool) {
	db.logger.Debug.Println("Getting user from database (by their Minecraft username)...")

	stmt, err := db.handler.FetchRowsFromTableStmt("users", []dbh.FilterItem{
		{"userName", dbh.CMP_EQ, nil},
	})

	defer stmt.Close()

	var uid, name, state string
	var userChange int

	if err != nil {
		goto errorState
	}

	err = stmt.QueryRow(username).Scan(&uid, &name, &state, &userChange)

errorState:
	if err != nil {

		return UserInfo{}, false

		db.logger.Debug.Println("Failed to obtain user from database. Invalid user(?)")
	} else {
		db.logger.Debug.Println("Managed to obtain user from database.")
	}

	return UserInfo{
		Uid:        uid,
		Name:       name,
		State:      state,
		UserChange: userChange,
	}, true

}

func (db *MCAuthDB_sqlite3) Getuser_changetime(id string) (skin_time int) {
	// start database transaction
	//tx, err := database.Begin()
	db.logger.Debug.Println("Getting user from database (by their Discord ID)...")

	stmt, err := db.handler.FetchRowsFromTableStmt("users", []dbh.FilterItem{
		{"discordID", dbh.CMP_EQ, nil},
	})

	defer stmt.Close()

	var uid, name, state string

	if err != nil {
		goto errorState
	}

	err = stmt.QueryRow(id).Scan(&uid, &name, &state, &skin_time)

	//fmt.Printf("test = %T\n", id)
errorState:
	if err != nil {
		// shit went fuck, probably an invalid username or something.
		skin_time = 0
		db.logger.Debug.Println("Failed to obtain user from database. Invalid user(?)")
	} else {
		db.logger.Debug.Println("Managed to obtain user from database.")
	}

	return

}
