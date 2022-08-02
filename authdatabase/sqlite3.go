package authdatabase

import (
	"database/sql"
	"log"
	"path/filepath"

	"git.dsrt-int.net/actionmc/actionmc-site-go/config"

	"git.dsrt-int.net/actionmc/actionmc-site-go/logging"
	dbh "git.dsrt-int.net/actionmc/actionmc-site-go/sqlite3dbh"
	_ "github.com/mattn/go-sqlite3"
)

type MCAuthDB_sqlite3 struct {
	handler *dbh.DBHandler
	logger  *logging.Logger
}

func createDB(handler *dbh.DBHandler) error {

	// TODO(ultrabear) some columns here use the 32 bit sql INT type to store unix timestamps
	// Anyone familiar with the 2038 problem will know that this will overflow in 2038
	// To solve this, all unix timestamps should be moved to BIGINT or T_INT64 and all databases should be migrated
	// This TODO is not immediate but will destroy the site if it is not patched

	var err error

	if _, err = handler.MakeNewTable("users", [][2]string{
		{"discordID", dbh.T_STR + dbh.T_PK},
		{"userName", dbh.T_STR},
		{"state", dbh.T_INT},
		// userChange is a unix timestamp
		{"userChange", dbh.T_INT}, // TODO(ultrabear) migrate databases to T_INT64 and move this to T_INT64
	}); err != nil {
		goto errorState
	}

	if _, err = handler.MakeNewTable("services", [][2]string{
		{"serviceIP", dbh.T_STR + dbh.T_PK},
		{"serviceName", dbh.T_STR},
		{"state", dbh.T_INT},
	}); err != nil {
		goto errorState
	}

	if _, err = handler.MakeNewTable("admins", [][2]string{
		{"discordID", dbh.T_STR + dbh.T_PK},
		{"privilegeLevel", dbh.T_INT},
	}); err != nil {
		goto errorState
	}

	if _, err = handler.MakeNewTable("deleted_users", [][2]string{
		{"shaSum", dbh.T_STR + dbh.T_PK},
		{"banned", dbh.T_INT},
	}); err != nil {
		goto errorState
	}

	if _, err = handler.MakeNewTable("reset_requests", [][2]string{
		{"discordID", dbh.T_STR + dbh.T_PK},
		// requestTime is a unix timestamp
		{"requestTime", dbh.T_INT}, // TODO(ultrabear) migrate databases to T_INT64 and move this to T_INT64
	}); err != nil {
		goto errorState
	}

	if _, err = handler.MakeNewTable("amgmt_tickets", [][2]string{
		{"discordID", dbh.T_STR + dbh.T_PK},
		{"dlTicket", dbh.T_STR},
		{"serverIp", dbh.T_STR},
	}); err != nil {
		goto errorState
	}

	if _, err = handler.MakeNewTable("sult_keys", [][2]string{
		{"discordID", dbh.T_STR + dbh.T_PK},
		{"totpSecret", dbh.T_STR},
	}); err != nil {
		goto errorState
	}

	if _, err = handler.MakeNewTable("cookie_consents", [][2]string{
		{"consentId", dbh.T_STR + dbh.T_PK},
		{"version", dbh.T_INT},
		// allow* are T_INTs being used as booleans
		{"allowNecessary", dbh.T_INT},
		{"allowPreferences", dbh.T_INT},
		// consentTimestamp is a unix timestamp
		{"consentTimestamp", dbh.T_INT64},
	}); err != nil {
		goto errorState
	}

	if _, err = handler.MakeNewTable("site_gallery", [][2]string{
		{"imageNumber", dbh.T_INT64 + dbh.T_PK},
		{"name", dbh.T_STR},
		{"description", dbh.T_STR},
		{"image_url", dbh.T_STR},
		// hasCredit is a T_INT used as boolean
		{"hasCredit", dbh.T_INT},
		{"creditedAuthor", dbh.T_STR},
		// publishTimestamp is a unix timestamp
		{"publishTimestamp", dbh.T_INT64},
	}); err != nil {
		goto errorState
	}

errorState:
	return err
}

type UserInfo struct {
	Uid        string
	Name       string
	State      string
	UserChange int
}

type UserError struct {
	err    string
	exists bool
	dberr  bool
}

func (UE UserError) Error() string    { return UE.err }
func (UE UserError) UserExists() bool { return UE.exists }
func (UE UserError) DBError() bool    { return UE.dberr }

// db initialiser
func InitSQLite3DB(logger *logging.Logger) *MCAuthDB_sqlite3 {
	diskdb, err := sql.Open("sqlite3", filepath.Join(config.GetDataDir(), "actionid.db"))
	logger.Debug.Println("Database opened, creating tables if they don't exist...")

	authDB := MCAuthDB_sqlite3{
		handler: dbh.OpenFromDB(diskdb),
		logger:  logger,
	}

	err = createDB(authDB.handler)

	if err != nil {
		logger.Err.Println("Error occurred while initialising on-disk database - more information below.")
		log.Printf("%q: %s\n", err)
	}

	logger.Debug.Println("Database initialised OK.")

	return &authDB
}
