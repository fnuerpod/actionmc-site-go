package sqlite3dbh

import "testing"

var dummydb = new(DBHandler)

const fmterr = "ERROR:\nGOT: %s\nWANT: %s"

func asserteq(t *testing.T, str, must interface{}) {

	// If it panics here who cares clearly the test input is busted
	// DO NOT type assume in production code, this is a test suite where it doesnt matter
	if s, oks := str.(string); oks {
		if s != must.(string) {
			t.Errorf(fmterr, str, must)
		}
	} else if n, okn := str.(int); okn {
		if n != must.(int) {
			t.Errorf(fmterr, str, must)
		}
	}
}

// MakeNewTable
func TestMakeNewTable(t *testing.T) {
	var str string

	str = dummydb.MakeNewTableBuilder("table", [][2]string{
		{"values", T_STR + T_PK},
		{"maps", T_INT},
	})

	asserteq(t, str, "CREATE TABLE IF NOT EXISTS 'table' (values TEXT PRIMARY KEY, maps INT)")

}

// AddToTable
func TestAddToTable(t *testing.T) {
	var str string

	str, _ = dummydb.AddToTableBuilder("table", []TableItem{
		{"value", 45},
		{"mapping", "dfg"},
	})

	asserteq(t, str, "REPLACE INTO 'table' (value, mapping)\nVALUES (?, ?)")

}

// FetchRowsFromTable
func TestFetchRowsFromTable(t *testing.T) {
	var str string

	str, _ = dummydb.FetchRowsFromTableBuilder("table", nil)

	asserteq(t, str, "SELECT * FROM 'table'")

	str, _ = dummydb.FetchRowsFromTableBuilder("table", []FilterItem{
		{"value", CMP_EQ, nil},
	})

	asserteq(t, str, "SELECT * FROM 'table' WHERE (value = ?)")

	str, _ = dummydb.FetchRowsFromTableBuilder("table", []FilterItem{
		{"value", CMP_EQ, nil},
		{"data", "LIKE", nil},
	})

	asserteq(t, str, "SELECT * FROM 'table' WHERE (value = ?) AND (data LIKE ?)")

}

// UpdateInTable
func TestUpdateInTable(t *testing.T) {
	var str string

	str, _ = dummydb.UpdateInTableBuilder("table", []TableItem{
		{"value", "45"},
	}, []FilterItem{
		{"map", CMP_EQ, 4},
	})

	asserteq(t, str, "UPDATE 'table' SET value = ? WHERE (map = ?)")

	str, _ = dummydb.UpdateInTableBuilder("table", []TableItem{
		{"value", "45"},
		{"string", 4},
	}, []FilterItem{
		{"map", CMP_EQ, 4},
		{"asd", CMP_NEQ, 34},
	})

	asserteq(t, str, "UPDATE 'table' SET value = ?, string = ? WHERE (map = ?) AND (asd != ?)")

	str, _ = dummydb.UpdateInTableBuilder("table", []TableItem{
		{"value", "45"},
		{"string", 4},
	}, nil)

	asserteq(t, str, "UPDATE 'table' SET value = ?, string = ?")

}

// DeleteRowsFromTable
func TestDeleteRowsFromTable(t *testing.T) {
	var str string

	// Test having a single conditional
	str, _ = dummydb.DeleteRowsFromTableBuilder("table", []FilterItem{
		{"value", CMP_EQ, nil},
	})

	asserteq(t, str, "DELETE FROM 'table' WHERE (value = ?)")

	// Test not having a conditional
	str, _ = dummydb.DeleteRowsFromTableBuilder("table", nil)

	asserteq(t, str, "DELETE FROM 'table'")

	// Test multi item
	str, _ = dummydb.DeleteRowsFromTableBuilder("table", []FilterItem{
		{"value", CMP_EQ, nil},
		{"map", CMP_EQ, nil},
	})

	asserteq(t, str, "DELETE FROM 'table' WHERE (value = ?) AND (map = ?)")

}

// DeleteTable
func TestDeleteTable(t *testing.T) {
	var str string

	// Test deleting a table
	str = dummydb.DeleteTableBuilder("table")

	asserteq(t, str, "DROP TABLE IF EXISTS 'table'")
}

// FetchTable
func TestFetchTable(t *testing.T) {
	var str string

	str = dummydb.FetchTableBuilder("table")

	asserteq(t, str, "SELECT * FROM 'table'")

}
