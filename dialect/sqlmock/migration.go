package sqlmock

import (
	m "github.com/silas/jdb/dialect/migration"
)

const version = 1

const createTable = `
CREATE TABLE
`

var revisions = m.Revisions{
	m.SQL(1, createTable),
}
