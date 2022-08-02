package authdatabase

import (
	//"bufio"
	//"bytes"
	"crypto/sha256"
	//"database/sql"
	"errors"
	//"git.dsrt-int.net/actionmc/actionmc-site-go/logging"
	dbh "git.dsrt-int.net/actionmc/actionmc-site-go/sqlite3dbh"
	_ "github.com/mattn/go-sqlite3"
	//"log"
	//"strconv"
)

func (db *MCAuthDB_sqlite3) Checkdeleted(rawid string) (is_banned bool) {
	db.logger.Debug.Println("Checking if user is banned (by SHA256'd Discord ID)...")

	hashbytes := sha256.Sum256([]byte(rawid))
	hash := string(hashbytes[:])

	stmt, err := db.handler.FetchRowsFromTableStmt("deleted_users", []dbh.FilterItem{
		{"shaSum", dbh.CMP_EQ, nil},
	})

	defer stmt.Close()

	var sumid string
	var isban bool

	if err != nil {
		goto errorState
	}

	err = stmt.QueryRow(hash).Scan(&sumid, &isban)

errorState:
	if err != nil {
		// shit went fuck, probably an invalid username or something.
		db.logger.Debug.Println("Failed to get user from database - is not banned.")
		is_banned = false
	} else {
		db.logger.Debug.Println("Managed to get user from database - is banned.")
		is_banned = true
	}

	return

}

func (db *MCAuthDB_sqlite3) Adddeleted_banned(id string) (bool, string) {
	// start transaction.
	db.logger.Debug.Println("Adding new banned user who deleted account...")

	_, err := db.handler.AddToTable("deleted_users", []dbh.TableItem{
		{"shaSum", sha256.Sum256([]byte(id))},
		{"banned", true},
	})

	if err != nil {
		db.logger.Debug.Println("Failed to add user to deleted banned list (Exec error).")
		return false, "bandel_trans_prep_error"
	}

	db.logger.Debug.Println("User added to banned deleted list successfully.")
	return true, ""
}

func (db *MCAuthDB_sqlite3) Deleteuser(id string) error {
	db.logger.Debug.Println("Deleting a user...")

	tx, err := db.handler.Begin()
	if err != nil {
		return errors.New("del_trans_init_error")
	}

	var txerr error

	for _, v := range []string{"users", "admins", "authenticator"} {

		// Operation on db.handler even though it operates on both dbs, this is because builder does not interact with the db
		str, args := db.handler.DeleteRowsFromTableBuilder(v, []dbh.FilterItem{
			{"discordID", dbh.CMP_EQ, id},
		})

		_, err := tx.Exec(str, args...)

		if err != nil {
			txerr = err
			goto errorState
		}

	}

errorState:
	if txerr != nil {
		// Failure, rollback and return error
		if err := tx.Rollback(); err != nil {
			db.logger.Err.Println("rollback failed")
		}
		return errors.New("rollback_error")
	} else {
		// Sucess, but will it commit
		if err := tx.Commit(); err != nil {
			db.logger.Err.Println("Error commiting data")
			return errors.New("commit_error")
		}
		db.logger.Debug.Println("deleted user")
		return nil
	}
}
