package authdatabase

import (
	"errors"

	dbh "git.dsrt-int.net/actionmc/actionmc-site-go/sqlite3dbh"
	_ "github.com/mattn/go-sqlite3"
)

func (db *MCAuthDB_sqlite3) GetGDPR(hash string) (exists bool, version int, allow_necessary bool, allow_preferences bool, time int64) {
	stmt, err := db.handler.FetchRowsFromTableStmt("cookie_consents", []dbh.FilterItem{
		{"consentId", dbh.CMP_EQ, nil},
	})

	defer stmt.Close()

	var allow_nec int
	var allow_pre int
	var ver int

	var ts int64
	var hasha string

	if err != nil {
		// shit went fuck, probably an invalid username or something.
		db.logger.Debug.Println("Failed to get user from database - no cookie consent.")
		exists = false
		version = 0
		allow_necessary = false
		allow_preferences = false
		ts = 0
	}

	err = stmt.QueryRow(hash).Scan(&hasha, &ver, &allow_nec, &allow_pre, &ts)

	if err != nil {
		// shit went fuck, probably an invalid username or something.
		db.logger.Debug.Println("Failed to get user from database - no cookie consent.")
		exists = false
		version = 0
		allow_necessary = false
		allow_preferences = false
		ts = 0
	} else {
		exists = true
		version = ver
		allow_necessary = (allow_nec == 1)
		allow_preferences = (allow_pre == 1)
		time = ts
	}

	return
}

func (db *MCAuthDB_sqlite3) Deleteconsent(id string) error {
	db.logger.Debug.Println("Deleting a GDPR consent...")

	tx, err := db.handler.Begin()
	if err != nil {
		return errors.New("del_trans_init_error")
	}

	var txerr error

	str, args := db.handler.DeleteRowsFromTableBuilder("cookie_consents", []dbh.FilterItem{
		{"consentId", dbh.CMP_EQ, id},
	})

	_, err = tx.Exec(str, args...)

	if err != nil {
		txerr = err
		goto errorState
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
		db.logger.Debug.Println("deleted GDPR consent")
		return nil
	}
}

func (db *MCAuthDB_sqlite3) AddGDPR(consentID string, version int, allow_necessary bool, allow_preferences bool, timestamp int64) error {
	// start transaction.
	db.logger.Debug.Println("Adding new GDPR consent...")

	_, err := db.handler.AddToTable("cookie_consents", []dbh.TableItem{
		{"consentId", consentID},
		{"version", version},
		{"allowNecessary", allow_necessary},
		{"allowPreferences", allow_preferences},
		{"consentTimestamp", timestamp},
	})

	if err != nil {
		db.logger.Fatal.Fatalln(err)
	}

	db.logger.Debug.Println("Consent created successfully.")
	return nil
}
