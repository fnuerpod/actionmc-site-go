package sqlite3dbh

import (
	"database/sql"
)

// DBHandler is a wrapper around a sql.DB that provides functions to not write pure sql
type DBHandler struct {
	*sql.DB
}

// Datatypes

// TableItem represents a piece of data to be sent to a sql column.
// It is written so that TableItem{"name", "tom"} is valid and readable
type TableItem struct {
	// SQL UNSAFE
	ColName string
	// SQL SAFE
	Data interface{}
}

// FilterItem represents a sql conditional statement in a more syntax safe way.
// It is written so that FilterItem{"name", "=", "tom"} is human readable and valid
type FilterItem struct {
	// SQL UNSAFE
	ColName string
	// SQL UNSAFE
	Comparator string
	// SQL SAFE
	Data interface{}
}

// Common sql col types

// Predefined int constants
const (
	T_INT8  = "TINYINT"
	T_INT16 = "SMALLINT"
	T_INT24 = "MEDIUMINT"
	T_INT32 = "INT"
	T_INT64 = "BIGINT"
	T_INT   = "INT"
)

// Predefined text constants
const (
	T_STR8  = "TINYTEXT"
	T_STR16 = "TEXT"
	T_STR24 = "MEDIUMTEXT"
	T_STR32 = "LONGTEXT"
	T_STR   = "TEXT"
)

// Predefined blob constants
const (
	T_BYTE8  = "TINYBLOB"
	T_BYTE16 = "BLOB"
	T_BYTE24 = "MEDIUMBLOB"
	T_BYTE32 = "LONGBLOB"
	T_BYTE   = "BLOB"
)

// Extras
const (
	T_PK  = " PRIMARY KEY"
	T_NIL = "NULL"
)

// Comparator types
const (
	CMP_EQ   = "="
	CMP_NEQ  = "!="
	CMP_GT   = ">"
	CMP_LT   = "<"
	CMP_GTEQ = ">="
	CMP_LTEQ = "<="
)
