package authdatabase

import (
	//"bufio"
	//"bytes"
	//"crypto/sha256"
	"database/sql"
	//"errors"
	//"git.dsrt-int.net/actionmc/actionmc-site-go/logging"
	"log"
	"strconv"

	dbh "git.dsrt-int.net/actionmc/actionmc-site-go/sqlite3dbh"
	_ "github.com/mattn/go-sqlite3"
)

func (db *MCAuthDB_sqlite3) Createuser(username string, id string) error {
	// start transaction.
	db.logger.Debug.Println("Registering new user...")

	// first check if this username exists.
	_, ok := db.Getstate_byuname(username)

	if ok {
		// user exists
		return &UserError{"user_exists_create", true, false}
	}

	_, err := db.handler.AddToTable("users", []dbh.TableItem{
		{"discordID", id},
		{"userName", username},
		{"state", 1},
		{"userChange", 0},
	})

	if err != nil {
		db.logger.Err.Println("Failed to update existing user (transaction execution error).")
		return &UserError{"user_trans_exec_error", false, true}
	}

	db.logger.Debug.Println("User created successfully.")
	return nil
}

func (db *MCAuthDB_sqlite3) CreateTOTP(id string, secret string) error {
	// start transaction.
	db.logger.Debug.Println("Registering new TOTP secret...")

	// first check if this username exists.
	_, ok := db.GetTOTPSecret(id)

	if ok {
		// user exists
		return &UserError{"totp_exists", true, false}
	}

	_, err := db.handler.AddToTable("sult_keys", []dbh.TableItem{
		{"discordID", id},
		{"totpSecret", secret},
	})

	if err != nil {
		db.logger.Err.Println("Failed to add TOTP key (transaction execution error).")
		return &UserError{"totp_trans_exec_error", false, true}
	}

	db.logger.Debug.Println("TOTP secret saved successfully.")
	return nil
}

// Changes a users username
// Holy fuck funey this name makes no sense pleae comment important shit like this
func (db *MCAuthDB_sqlite3) Updateuser(username string, id string) error {
	// start transaction.
	db.logger.Debug.Println("Updating existing user...")

	// first check if this username exists.
	_, ok := db.Getstate_byuname(username)

	if ok {
		// username is taken
		return &UserError{"user_exists_update", true, false}
	}

	_, err := db.handler.UpdateInTable("users", []dbh.TableItem{
		{"userName", username},
	}, []dbh.FilterItem{
		{"discordID", dbh.CMP_EQ, id},
	})

	if err != nil {
		db.logger.Err.Println("Failed to update existing user (transaction execution error).")
		return &UserError{"user_trans_exec_error", false, true}
	}

	db.logger.Debug.Println("User updated successfully.")
	return nil
}

func (db *MCAuthDB_sqlite3) Updateuser_skintime(id string, skin_change int) error {
	// start transaction.
	db.logger.Debug.Println("Updating existing user skin change time...")

	_, err := db.handler.UpdateInTable("users", []dbh.TableItem{
		{"userChange", skin_change},
	}, []dbh.FilterItem{
		{"discordID", dbh.CMP_EQ, id},
	})

	if err != nil {
		db.logger.Err.Println("Failed to update existing user's time (transaction execution error).")
		return &UserError{"user_time_trans_exec_error", false, true}
	}

	db.logger.Debug.Println("User updated time successfully.")
	return nil
}

func (db *MCAuthDB_sqlite3) Getall() (rows *sql.Rows) {
	db.logger.Debug.Println("Getting all users from database...")
	rows, err := db.handler.FetchTable("users")
	if err != nil {
		log.Fatal(err)
	}
	db.logger.Debug.Println("Got all users from database.")
	return
}

func (db *MCAuthDB_sqlite3) Getallservices() (rows *sql.Rows) {
	db.logger.Debug.Println("Getting all services from database...")
	rows, err := db.handler.FetchTable("services")
	if err != nil {
		log.Fatal(err)
	}
	db.logger.Debug.Println("Got all services from database.")
	return
}

// THIS JUST FUCKING PANICS IT DOESNT NEED TO RETURN ANYTHING I SWEAR TO FUCK FUNEY
func (db *MCAuthDB_sqlite3) Changestate(id string, state string) {
	db.logger.Debug.Println("Changing state of user...")

	istate, err := strconv.Atoi(state)

	if err != nil {
		log.Panic(err)
	}

	_, err = db.handler.UpdateInTable("users", []dbh.TableItem{
		{"state", istate},
	}, []dbh.FilterItem{
		{"discordID", dbh.CMP_EQ, id},
	})

	if err != nil {
		db.logger.Fatal.Fatalln("Failed to update user, this is fatal - program will terminate. More information below...")
		log.Panic(err)
	}

	db.logger.Debug.Println("Changed state of user.")
}
