package jdb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestDeleteBuilder_ToSQL(t *testing.T) {
	c, mock := createMockClient(t)
	defer c.Close()

	kind := "test"
	s, p, err := c.Query(kind).Delete().ToSQL()
	require.NoError(t, err)
	require.Equal(t, "DELETE FROM jdb WHERE ((kind = ?))", s)
	require.Equal(t, params(kind), p)

	s, p, err = c.Query(kind).Delete("1").ToSQL()
	require.NoError(t, err)
	require.Equal(t, "DELETE FROM jdb WHERE ((kind = ?) AND (id = ?))", s)
	require.Equal(t, params(kind, "1"), p)

	s, p, err = c.Query(kind).Delete("1", "2").ToSQL()
	require.NoError(t, err)
	require.Equal(t, "DELETE FROM jdb WHERE ((kind = ?) AND (id IN (?, ?)))", s)
	require.Equal(t, params(kind, "1", "2"), p)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteBuilder_Exec(t *testing.T) {
	c, mock := createMockClient(t)
	defer c.Close()

	kind := "test"

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM jdb WHERE").
		WithArgs(kind).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("DELETE FROM jdb WHERE").
		WithArgs(kind, "1", "2").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	c.Tx(context.Background(), func(tx *Tx) error {
		err := c.Query(kind).Delete().Exec(context.Background(), tx)
		require.NoError(t, err)

		err = c.Query(kind).Delete("1", "2").Exec(context.Background(), tx)
		require.NoError(t, err)

		return tx.Commit()
	})

	require.NoError(t, mock.ExpectationsWereMet())
}
