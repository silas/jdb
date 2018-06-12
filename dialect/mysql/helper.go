package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/silas/jdb/dialect/migration"
)

const getID = `
SELECT
  json_extract(data, '$.ID')
FROM
  {{ .Table }}
WHERE
  kind = 'jdb' AND
  id = ?
LIMIT 1;
`

const setID = `
INSERT INTO {{ .Table }}
  (kind, id, data)
VALUES
  ('jdb', ?, ?)
ON DUPLICATE KEY UPDATE
  data = ?,
  update_time = ` + timestamp + `;
`

const tableExists = `
SELECT
  count(*) > 0 AS table_exists
FROM
  information_schema.tables
WHERE
  table_schema = database() AND
  table_name = ?
LIMIT 1;
`

type migrationHelper struct {
	table string
}

func (h *migrationHelper) Version() int {
	return version
}

func (h *migrationHelper) getID(ctx context.Context, tx *sql.Tx, name string) (int, error) {
	var id string
	err := tx.QueryRowContext(ctx, h.RenderSQL(getID), name).Scan(&id)
	if err != nil {
		return 0, err
	}

	version, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}
	return version, nil
}

func (h *migrationHelper) setID(ctx context.Context, tx *sql.Tx, name string, id int) error {
	data := fmt.Sprintf(`{"ID":%d}`, id)
	_, err := tx.ExecContext(ctx, h.RenderSQL(setID), name, data, data)
	return err
}

func (h *migrationHelper) GetRevision(ctx context.Context, tx *sql.Tx) (int, error) {
	return h.getID(ctx, tx, "revision")
}

func (h *migrationHelper) SetRevision(ctx context.Context, tx *sql.Tx, id int) error {
	return h.setID(ctx, tx, "revision", id)
}

func (h *migrationHelper) GetVersion(ctx context.Context, tx *sql.Tx) (int, error) {
	return h.getID(ctx, tx, "version")
}

func (h *migrationHelper) SetVersion(ctx context.Context, tx *sql.Tx, id int) error {
	return h.setID(ctx, tx, "version", id)
}

func (h *migrationHelper) TableExists(ctx context.Context, tx *sql.Tx) (bool, error) {
	var exists bool
	err := tx.QueryRowContext(ctx, h.RenderSQL(tableExists), h.table).Scan(&exists)
	return exists, err
}

func (h *migrationHelper) RenderSQL(text string, locals ...map[string]interface{}) string {
	data := []map[string]interface{}{
		{
			"Table":     h.table,
			"Namespace": h.table,
		},
	}
	data = append(data, locals...)
	return migration.Render("sql", text, data...)
}
