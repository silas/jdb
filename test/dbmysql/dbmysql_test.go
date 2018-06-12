package dbmysql

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/silas/jdb/dialect/mysql"
	"github.com/silas/jdb/test/db"
	"github.com/stretchr/testify/require"
)

func TestMySQL(t *testing.T) {
	driverName := "mysql"
	dataSourceName := os.Getenv("JDB_MYSQL_DSN")
	table := "jdb_test"

	if testing.Short() || dataSourceName == "" {
		t.Skip("skipping mysql testdb in short mode")
	}

	c, err := sql.Open(driverName, dataSourceName)
	require.NoError(t, err)
	defer c.Close()

	_, err = c.Exec("DROP TABLE IF EXISTS " + table)
	require.NoError(t, err)
	c.Close()

	db.New(driverName, dataSourceName, dataSourceName, table).Run(t)
}
