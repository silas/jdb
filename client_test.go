package jdb

import (
	"context"
	"fmt"
	"testing"
	"time"

	jdbsqlmock "github.com/silas/jdb/dialect/sqlmock"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestOpen(t *testing.T) {
	c, mock := createMockClient(t)
	defer c.Close()

	require.NotNil(t, c.ID)
	require.NotNil(t, c.Kind)
	require.NotNil(t, c.ParentKind)
	require.NotNil(t, c.ParentId)
	require.NotNil(t, c.UniqueStringKey)
	require.NotNil(t, c.StringKey)
	require.NotNil(t, c.NumericKey)
	require.NotNil(t, c.TimeKey)
	require.NotNil(t, c.Data)
	require.NotNil(t, c.CreateTime)
	require.NotNil(t, c.UpdateTime)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOpen_Option_ReadOnly(t *testing.T) {
	_, err := Open("sqlmock", "require-read-only-true", ReadOnly(true))
	require.NoError(t, err)
	_, err = Open("sqlmock", "require-read-only-false", ReadOnly(true))
	require.EqualError(t, err, "ReadOnly=true")

	_, err = Open("sqlmock", "require-read-only-false", ReadOnly(false))
	require.NoError(t, err)
	_, err = Open("sqlmock", "require-read-only-false")
	require.NoError(t, err)
	_, err = Open("sqlmock", "require-read-only-true", ReadOnly(false))
	require.EqualError(t, err, "ReadOnly=false")
}

func TestOpen_Option_Table(t *testing.T) {
	c, err := Open("sqlmock", "test")
	require.NoError(t, err)
	require.Equal(t, c.table, "jdb")

	c, err = Open("sqlmock", "test", Table("test"))
	require.NoError(t, err)
	require.Equal(t, c.table, "test")

	_, err = Open("sqlmock", "test", Table(""))
	require.EqualError(t, err, "jdb: invalid table name")
}

func TestClient_Migrate(t *testing.T) {
	c, mock := createMockClient(t)
	defer c.Close()

	jdbsqlmock.ExpectRun(mock)

	require.NoError(t, c.Migrate(context.Background()))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestClient_Path(t *testing.T) {
	c, mock := createMockClient(t)
	defer c.Close()

	require.Equal(t, c.Path("one").toWhereField(), "data->'$.one'")
	require.Equal(t, c.Path("one", "two").toWhereField(), "data->'$.one.two'")
	require.Equal(t, c.Path("one").Key("two").toWhereField(), "data->'$.one.two'")
	require.Equal(t, c.Path("one").Index(2).toWhereField(), "data->'$.one[2]'")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestClient_Tx(t *testing.T) {
	c, mock := createMockClient(t)
	defer c.Close()

	mock.ExpectBegin()
	mock.ExpectCommit()
	c.Tx(context.Background(), func(tx *Tx) error {
		return tx.Commit()
	})
	require.NoError(t, mock.ExpectationsWereMet())

	mock.ExpectBegin()
	mock.ExpectRollback()
	c.Tx(context.Background(), func(tx *Tx) error {
		return nil
	})
	require.NoError(t, mock.ExpectationsWereMet())
}

func init() {
	RegisterDialect(jdbsqlmock.RegisterDialectArgs())
}

func createMockClient(t *testing.T) (*Client, sqlmock.Sqlmock) {
	dsn := fmt.Sprintf("dsn-%d", time.Now().UnixNano())
	_, mock, err := sqlmock.NewWithDSN(dsn)
	require.NoError(t, err)

	c, err := Open("sqlmock", dsn)
	require.NoError(t, err)
	return c, mock
}
