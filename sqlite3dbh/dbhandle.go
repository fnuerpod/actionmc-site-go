package sqlite3dbh

import (
	"database/sql"
)

// MakeNewTable

// MakeNewTable creates a new table in the database with the given name/type pairing
func (db *DBHandler) MakeNewTable(tableName string, cols [][2]string) (sql.Result, error) {
	return db.Exec(db.MakeNewTableBuilder(tableName, cols))
}

// MakeNewTableStmt returns a sql.Stmt of MakeNewTable on the DB object
func (db *DBHandler) MakeNewTableStmt(tableName string, cols [][2]string) (*sql.Stmt, error) {
	return db.Prepare(db.MakeNewTableBuilder(tableName, cols))
}

// AddToTable

// AddToTable adds a row to the given table
func (db *DBHandler) AddToTable(tableName string, cols []TableItem) (sql.Result, error) {
	str, args := db.AddToTableBuilder(tableName, cols)
	return db.Exec(str, args...)
}

// AddToTableStmt returns a sql.Stmt of AddToTable on the DB object, TableItem.Data is ignored
func (db *DBHandler) AddToTableStmt(tableName string, cols []TableItem) (*sql.Stmt, error) {
	str, _ := db.AddToTableBuilder(tableName, cols)
	return db.Prepare(str)
}

// FetchRowsFromTable

// FetchRowsFromTable fetches rows from a table that match the search criteria
func (db *DBHandler) FetchRowsFromTable(tableName string, filters []FilterItem) (*sql.Rows, error) {
	str, args := db.FetchRowsFromTableBuilder(tableName, filters)
	return db.Query(str, args...)
}

// FetchRowsFromTableStmt returns a sql.Stmt of FetchRowsFromTable on the DB object, FilterItem.Data is ignored
func (db *DBHandler) FetchRowsFromTableStmt(tableName string, filters []FilterItem) (*sql.Stmt, error) {
	str, _ := db.FetchRowsFromTableBuilder(tableName, filters)
	return db.Prepare(str)
}

// FetchTable

// FetchTable fetches the entire table
func (db *DBHandler) FetchTable(tableName string) (*sql.Rows, error) {
	return db.Query(db.FetchTableBuilder(tableName))
}

// FetchTableStmt returns a sql.Stmt of FetchTable on the DB object
func (db *DBHandler) FetchTableStmt(tableName string) (*sql.Stmt, error) {
	return db.Prepare(db.FetchTableBuilder(tableName))
}

// UpdateInTable

// UpdateInTable updates rows in a table with the given filters (or nil filters for update all)
func (db *DBHandler) UpdateInTable(tableName string, edits []TableItem, filters []FilterItem) (sql.Result, error) {
	str, args := db.UpdateInTableBuilder(tableName, edits, filters)
	return db.Exec(str, args...)
}

// UpdateInTableStmt returns a sql.Stmt of UpdateInTable on the DB object, (Filter|Table)Item.Data is ignored
func (db *DBHandler) UpdateInTableStmt(tableName string, edits []TableItem, filters []FilterItem) (*sql.Stmt, error) {
	str, _ := db.UpdateInTableBuilder(tableName, edits, filters)
	return db.Prepare(str)
}

// DeleteRowsFromTable

// DeleteRowsFromTable deletes the rows that match the filter struct(s)
func (db *DBHandler) DeleteRowsFromTable(tableName string, filters []FilterItem) (sql.Result, error) {
	str, args := db.DeleteRowsFromTableBuilder(tableName, filters)
	return db.Exec(str, args...)
}

// DeleteRowsFromTableStmt returns a sql.Stmt of DeleteRowsFromTable on the DB object
func (db *DBHandler) DeleteRowsFromTableStmt(tableName string, filters []FilterItem) (*sql.Stmt, error) {
	str, _ := db.DeleteRowsFromTableBuilder(tableName, filters)
	return db.Prepare(str)
}

// DeleteTable

// DeleteTable deletes the given table
func (db *DBHandler) DeleteTable(tableName string) (sql.Result, error) {
	return db.Exec(db.DeleteTableBuilder(tableName))
}

// DeleteTableStmt returns a sql.Stmt of DeleteTable on the DB object
func (db *DBHandler) DeleteTableStmt(tableName string) (*sql.Stmt, error) {
	return db.Prepare(db.DeleteTableBuilder(tableName))
}
