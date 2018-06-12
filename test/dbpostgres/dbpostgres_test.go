package dbpostgres

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/silas/jdb/dialect/postgres"
	"github.com/silas/jdb/test/db"
	"github.com/stretchr/testify/require"
)

func TestPostgres(t *testing.T) {
	driverName := "postgres"
	dataSourceName := os.Getenv("JDB_POSTGRES_DSN")
	table := "jdb_test"

	if testing.Short() || dataSourceName == "" {
		t.Skip("skipping postgres testdb in short mode")
	}

	c, err := sql.Open(driverName, dataSourceName)
	require.NoError(t, err)
	defer c.Close()

	_, err = c.Exec("DROP TABLE IF EXISTS " + table)
	require.NoError(t, err)
	c.Close()

	db.New(driverName, dataSourceName, dataSourceName, table).Run(t)
}
