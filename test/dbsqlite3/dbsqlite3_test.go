package dbsqlite3

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/silas/jdb/dialect/sqlite3"
	"github.com/silas/jdb/test/db"
	"github.com/stretchr/testify/require"
)

func TestSqlite3(t *testing.T) {
	dir, err := ioutil.TempDir("", "jdb_test_")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "test.db")

	dataSourceName := path + "?cache=shared"
	readOnlyDataSourceName := dataSourceName + "&mode=ro"

	db.New("sqlite3", dataSourceName, readOnlyDataSourceName, "jdb_test").Run(t)
}
