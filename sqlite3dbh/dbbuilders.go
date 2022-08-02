package sqlite3dbh

import (
	"strings"
)

// MakeNewTableBuilder is the builder for MakeNewTable functions
func (db *DBHandler) MakeNewTableBuilder(tableName string, cols [][2]string) string {

	const header = "CREATE TABLE IF NOT EXISTS '"

	var blen int

	// Add pre loop lengths
	blen += len(header) + len(tableName) + 3

	// Add length of each item
	for i := range cols {
		blen += len(cols[i][0]) + len(cols[i][1]) + 1
	}

	// Add for sep
	blen += (len(cols) - 1) * 2

	// Add last )
	blen++

	var b strings.Builder
	b.Grow(blen)

	// Write create tablename
	b.WriteString(header)
	b.WriteString(tableName)
	b.WriteString("' (")

	// Write each item similar to strings.Join
	for i := range cols {

		if i != 0 {
			b.WriteString(", ")
		}

		b.WriteString(cols[i][0])
		b.WriteByte(' ')
		b.WriteString(cols[i][1])

	}

	b.WriteByte(')')

	// Ram wastage check
	if b.Len() != blen {
		panic(b.Len() - blen)
	}

	return b.String()
}

// AddToTableBuilder is the builder for AddToTable functions
func (db *DBHandler) AddToTableBuilder(tableName string, cols []TableItem) (string, []interface{}) {

	const header = "REPLACE INTO '"

	var blen int

	// Header length
	blen += len(header) + len(tableName) + 3

	// names sep
	blen += (len(cols) - 1) * 2

	// names items
	for i := range cols {
		blen += len(cols[i].ColName)
	}

	// VALUE string
	blen += len(")\nVALUES (")

	// sep and items of placeholder string and last )
	blen += (len(cols)-1)*2 + len(cols) + 1

	var b strings.Builder
	b.Grow(blen)

	b.WriteString(header)
	b.WriteString(tableName)
	b.WriteString("' (")

	// Write col names
	for i := range cols {
		if i != 0 {
			b.WriteString(", ")
		}

		b.WriteString(cols[i].ColName)
	}

	b.WriteString(")\nVALUES (")

	// Make data args list
	args := make([]interface{}, 0, len(cols))

	// Write data placeholders
	for i := range cols {
		if i != 0 {
			b.WriteString(", ")
		}

		b.WriteByte('?')

		// Write arg
		args = append(args, cols[i].Data)
	}

	b.WriteByte(')')

	// Ram wastage check
	if b.Len() != blen {
		panic(b.Len() - blen)
	}

	return b.String(), args
}

// FetchRowsFromTableBuilder is the builder for FetchRowsFromTable functions
func (db *DBHandler) FetchRowsFromTableBuilder(tableName string, filters []FilterItem) (string, []interface{}) {

	const header = "SELECT * FROM '"

	var blen int

	// Add header lengths
	blen += len(header) + len(tableName)

	// Only add extra lens if filters is not nil/empty
	if len(filters) > 0 {
		blen += len("' WHERE ")

		// Add sep lengths
		blen += (len(filters) - 1) * 5

		// Add each item length
		for i := range filters {
			blen += len(filters[i].ColName) + len(filters[i].Comparator) + 5
		}
	} else {
		blen += 1
	}

	var b strings.Builder
	b.Grow(blen)

	// Add headers
	b.WriteString(header)
	b.WriteString(tableName)

	// Make data args list
	args := make([]interface{}, 0, len(filters))

	// allow for nil filters to act as a FetchTable
	if len(filters) > 0 {
		b.WriteString("' WHERE ")

		// Add each filter
		for i := range filters {
			if i != 0 {
				b.WriteString(" AND ")
			}

			b.WriteByte('(')
			b.WriteString(filters[i].ColName)
			b.WriteByte(' ')
			b.WriteString(filters[i].Comparator)
			b.WriteString(" ?)")

			// Add arg
			args = append(args, filters[i].Data)
		}
	} else {
		b.WriteByte('\'')
	}

	// Ram wastage check
	if b.Len() != blen {
		panic(b.Len() - blen)
	}

	return b.String(), args
}

// FetchTableBuilder is the builder for FetchTable functions, this is eqivalent to FetchRowsFromTable(tableName, nil)
func (db *DBHandler) FetchTableBuilder(tableName string) string {
	// Because we only add strings once, concat is faster than a builder
	return "SELECT * FROM '" + tableName + "'"
}

// UpdateInTableBuilder is the builder for UpdateInTable functions
func (db *DBHandler) UpdateInTableBuilder(tableName string, edits []TableItem, filters []FilterItem) (string, []interface{}) {

	const (
		header   = "UPDATE '"
		setstr   = "' SET "
		wherestr = " WHERE "
	)

	var blen int

	// Add headers len
	blen += len(header) + len(setstr) + len(tableName)

	// Add edits length
	for i := range edits {
		blen += len(edits[i].ColName) + 4
	}

	// Add sep length of edits
	blen += (len(edits) - 1) * 2

	if filters != nil {
		// Add wherestr
		blen += len(wherestr)

		// Add len of seps
		blen += (len(filters) - 1) * 5

		// Add each datas length
		for i := range filters {
			blen += len(filters[i].ColName) + len(filters[i].Comparator) + 5
		}
	}

	var b strings.Builder
	b.Grow(blen)

	b.WriteString(header)
	b.WriteString(tableName)
	b.WriteString(setstr)

	args := make([]interface{}, 0, len(edits)+len(filters))

	for i := range edits {
		if i != 0 {
			b.WriteString(", ")
		}

		b.WriteString(edits[i].ColName)
		b.WriteString(" = ?")

		args = append(args, edits[i].Data)
	}

	if filters != nil {
		b.WriteString(wherestr)

		for i := range filters {
			if i != 0 {
				b.WriteString(" AND ")
			}

			b.WriteByte('(')
			b.WriteString(filters[i].ColName)
			b.WriteByte(' ')
			b.WriteString(filters[i].Comparator)
			b.WriteString(" ?)")

			args = append(args, filters[i].Data)
		}

	}

	// Ram wastage check
	if b.Len() != blen {
		panic(b.Len() - blen)
	}

	return b.String(), args
}

// DeleteRowsFromTableBuilder is the builder for DeleteRowsFromTable functions
func (db *DBHandler) DeleteRowsFromTableBuilder(tableName string, filters []FilterItem) (string, []interface{}) {

	const (
		header   = "DELETE FROM '"
		wherestr = "' WHERE "
		andstr   = " AND "
	)

	var blen int

	// Add header lengths
	blen += len(header) + len(tableName)

	if len(filters) > 0 {
		blen += len(wherestr)

		// Add sep lengths
		blen += (len(filters) - 1) * len(andstr)

		// Add length of comparator+col
		for i := range filters {
			blen += len(filters[i].ColName) + len(filters[i].Comparator) + 5
		}
	} else {
		// len for '
		blen++
	}

	var b strings.Builder
	b.Grow(blen)

	b.WriteString(header)
	b.WriteString(tableName)

	if len(filters) > 0 {
		b.WriteString(wherestr)
	} else {
		b.WriteByte('\'')
	}

	args := make([]interface{}, 0, len(filters))

	for i := range filters {
		if i != 0 {
			b.WriteString(andstr)
		}

		b.WriteByte('(')
		b.WriteString(filters[i].ColName)
		b.WriteByte(' ')
		b.WriteString(filters[i].Comparator)
		b.WriteString(" ?)")

		args = append(args, filters[i].Data)
	}

	// Ram wastage check
	if b.Len() != blen {
		panic(b.Len() - blen)
	}

	return b.String(), args
}

// DeleteTableBuilder is the builder for DeleteTable functions
func (db *DBHandler) DeleteTableBuilder(tableName string) string {
	// Because we only add strings once it is faster to use concat than a builder
	return "DROP TABLE IF EXISTS '" + tableName + "'"
}
