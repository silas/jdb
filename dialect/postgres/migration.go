package postgres

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
  numeric_key DOUBLE PRECISION,
  time_key TIMESTAMP WITH TIME ZONE,
  data JSONB,
  create_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT ` + timestamp + `,
  update_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT ` + timestamp + `,
  PRIMARY KEY (kind, id),
  FOREIGN KEY (parent_kind, parent_id) REFERENCES {{ .Table }} (kind, id)
);
`

var revisions = m.Revisions{
	m.SQL(1, createTable),
	m.SQL(2, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (create_time NULLS FIRST);`),
	m.SQL(3, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (update_time NULLS FIRST);`),
	m.SQL(4, `CREATE UNIQUE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (kind NULLS FIRST, unique_string_key NULLS FIRST);`),
	m.SQL(5, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (kind NULLS FIRST, string_key NULLS FIRST);`),
	m.SQL(6, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (kind NULLS FIRST, numeric_key NULLS FIRST);`),
	m.SQL(7, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (kind NULLS FIRST, time_key NULLS FIRST);`),
	m.SQL(8, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (kind NULLS FIRST, create_time NULLS FIRST);`),
	m.SQL(9, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (kind NULLS FIRST, update_time NULLS FIRST);`),
	m.SQL(10, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (kind NULLS FIRST, id NULLS FIRST, parent_kind NULLS FIRST);`),
	m.SQL(11, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (parent_kind NULLS FIRST, parent_id NULLS FIRST, kind NULLS FIRST, unique_string_key NULLS FIRST);`),
	m.SQL(12, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (parent_kind NULLS FIRST, parent_id NULLS FIRST, kind NULLS FIRST, string_key NULLS FIRST);`),
	m.SQL(13, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (parent_kind NULLS FIRST, parent_id NULLS FIRST, kind NULLS FIRST, numeric_key NULLS FIRST);`),
	m.SQL(14, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (parent_kind NULLS FIRST, parent_id NULLS FIRST, kind NULLS FIRST, time_key NULLS FIRST);`),
	m.SQL(15, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (parent_kind NULLS FIRST, parent_id NULLS FIRST, kind NULLS FIRST, create_time NULLS FIRST);`),
	m.SQL(16, `CREATE INDEX {{ .Namespace }}_r{{ .ID }} ON {{ .Table }} (parent_kind NULLS FIRST, parent_id NULLS FIRST, kind NULLS FIRST, update_time NULLS FIRST);`),
}
