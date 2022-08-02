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
	"log"
	"strconv"
)

func (db *MCAuthDB_sqlite3) Changeservicestate(serviceip string, state string) {
	db.logger.Debug.Println("Changing state of service...")

	istate, err := strconv.Atoi(state)

	if err != nil {
		log.Panic(err)
	}

	_, err = db.handler.UpdateInTable("services", []dbh.TableItem{
		{"state", istate},
	}, []dbh.FilterItem{
		{"serviceIP", dbh.CMP_EQ, serviceip},
	})

	if err != nil {
		db.logger.Fatal.Fatalln("Failed to update service (transaction execution), this is fatal - program will terminate. More information below...")
		log.Panic(err)
	}

	db.logger.Debug.Println("Changed state of service.")
}
