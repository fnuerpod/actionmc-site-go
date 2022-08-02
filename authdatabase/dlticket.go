package authdatabase

import (
	dbh "git.dsrt-int.net/actionmc/actionmc-site-go/sqlite3dbh"
	_ "github.com/mattn/go-sqlite3"
)

func (db *MCAuthDB_sqlite3) GetTicket(dlticket string) (valid bool, discordID string, ipHash string) {
	db.logger.Debug.Println("Getting Download Ticket from database...")

	rows, err := db.handler.FetchRowsFromTable("amgmt_tickets", []dbh.FilterItem{
		{"dlTicket", dbh.CMP_EQ, dlticket},
	})
	defer rows.Close()

	var db_downloadTicket string

	if err != nil {
		goto errorState
	}

	if rows.Next() {
		err = rows.Scan(&discordID, &db_downloadTicket, &ipHash)
	} else {
		valid = false
		return
	}

errorState:
	if err != nil {
		// shit went fuck, probably an invalid username or something.
		db.logger.Debug.Println("Failed to get ticket from database - invalid.")
		valid = false
	} else {
		db.logger.Debug.Println("Managed to get ticket from database - valid.")
		valid = true
	}

	return

}
