package db

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/silas/jdb"
	"github.com/silas/jdb/dialect"
	"github.com/silas/jdb/test/db/internal/data"
	"github.com/stretchr/testify/require"
)

type Test struct {
	db                     *jdb.Client
	rdb                    *jdb.Client
	d                      dialect.Dialect
	driverName             string
	dataSourceName         string
	readOnlyDataSourceName string
	table                  string
	sqlDB                  *sql.DB
}

func New(driverName, dataSourceName, readOnlyDataSourceName, table string) *Test {
	dt := &Test{
		driverName:             driverName,
		dataSourceName:         dataSourceName,
		readOnlyDataSourceName: readOnlyDataSourceName,
		table: table,
	}
	return dt
}

func (dt *Test) Run(t *testing.T) {
	d, err := jdb.Dialect(dt.driverName)
	require.NoError(t, err)

	tableOpt := jdb.Table(dt.table)

	db, err := jdb.Open(dt.driverName, dt.dataSourceName, tableOpt)
	require.NoError(t, err)
	require.NoError(t, db.Migrate(context.Background()))

	rdb, err := jdb.Open(dt.driverName, dt.readOnlyDataSourceName, tableOpt, jdb.ReadOnly(true))
	require.NoError(t, err)
	require.NoError(t, db.Migrate(context.Background()))

	sqlDB, err := sql.Open(dt.driverName, dt.dataSourceName)
	require.NoError(t, err)

	dt.d = d
	dt.db = db
	dt.rdb = rdb
	dt.sqlDB = sqlDB

	dt.testClient(t)
	dt.testDelete(t)
	dt.testSelect(t)
	dt.testInsert(t)
	dt.testUpdate(t)
}

func (dt *Test) setup(t *testing.T, populate bool) *jdb.Client {
	tx, err := dt.sqlDB.Begin()
	require.NoError(t, err)

	dt.deleteAll(t)

	if populate {
		for _, row := range data.UserRows {
			var parentKind, parentID, uniqueStringKey, stringKey, json *string
			var numericKey *float64
			var timeKey *time.Time
			if row.ParentKind != "" {
				parentID = &row.ParentKind
			}
			if row.ParentID != "" {
				parentID = &row.ParentID
			}
			if row.UniqueStringKey != "" {
				uniqueStringKey = &row.UniqueStringKey
			}
			if row.StringKey != "" {
				stringKey = &row.StringKey
			}
			if row.NumericKey != 0 {
				numericKey = &row.NumericKey
			}
			if !row.TimeKey.IsZero() {
				timeKey = &row.TimeKey
			}
			if row.Data != "" {
				json = &row.Data
			}

			dt.insertRaw(t, row.Kind, row.ID, parentKind, parentID, uniqueStringKey, stringKey, numericKey, timeKey,
				json, row.CreateTime, row.UpdateTime)
		}
	}

	err = tx.Commit()
	require.NoError(t, err)

	return dt.db
}

func (dt *Test) sql(s string) string {
	s = strings.Replace(s, "jdb_test", dt.table, -1)
	s = dt.d.ReplacePlaceHolders(s)
	return s
}

func (dt *Test) count(t *testing.T, kind string) int {
	s := `SELECT count(*) FROM jdb_test WHERE kind = ?`
	var count int
	err := dt.sqlDB.QueryRow(dt.sql(s), kind).Scan(&count)
	require.NoError(t, err)
	return count
}

func (dt *Test) deleteAll(t *testing.T) {
	s := `DELETE FROM jdb_test WHERE kind != ?`
	dt.exec(t, s, "jdb")
}

func (dt *Test) insertRaw(t *testing.T, args ...interface{}) {
	s := `
INSERT INTO jdb_test
  (kind, id, parent_kind, parent_id, unique_string_key, string_key, numeric_key, time_key, data, create_time, update_time)
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	dt.exec(t, s, args...)
}

func (dt *Test) exec(t *testing.T, s string, args ...interface{}) {
	_, err := dt.sqlDB.Exec(dt.sql(s), args...)
	require.NoError(t, err)
}
