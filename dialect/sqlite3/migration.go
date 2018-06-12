package sqlite3

import (
	m "github.com/silas/jdb/dialect/migration"
)

const version = 1

const createTable = `
CREATE TABLE {{ .Table }} (
  kind VARCHAR(64),
  id VARCHAR(64),
  parent_kind VARCHAR(64),
  parent_id VARCHAR(64),
  unique_string_key VARCHAR(255),
  string_key VARCHAR(255),
  numeric_key REAL,
  time_key DATETIME,
  data JSON,
  create_time DATETIME NOT NULL DEFAULT(` + timestamp + `),
  update_time DATETIME NOT NULL DEFAULT(` + timestamp + `),
  PRIMARY KEY (kind, id),
  FOREIGN KEY (parent_kind, parent_id) REFERENCES {{ .Table }} (kind, id)
);
`

var revisions = m.Revisions{
	m.SQL(1, createTable),
	m.SQL(2, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (create_time);`),
	m.SQL(3, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (update_time);`),
	m.SQL(4, `CREATE UNIQUE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (kind, unique_string_key);`),
	m.SQL(5, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (kind, string_key);`),
	m.SQL(6, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (kind, numeric_key);`),
	m.SQL(7, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (kind, time_key);`),
	m.SQL(8, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (kind, create_time);`),
	m.SQL(9, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (kind, update_time);`),
	m.SQL(10, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (kind, id, parent_kind);`),
	m.SQL(11, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (parent_kind, parent_id, kind, unique_string_key);`),
	m.SQL(12, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (parent_kind, parent_id, kind, string_key);`),
	m.SQL(13, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (parent_kind, parent_id, kind, numeric_key);`),
	m.SQL(14, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (parent_kind, parent_id, kind, time_key);`),
	m.SQL(15, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (parent_kind, parent_id, kind, create_time);`),
	m.SQL(16, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (parent_kind, parent_id, kind, update_time);`),
}
