package jdb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestTx_Now(t *testing.T) {
	c, mock := createMockClient(t)
	defer c.Close()

	now := time.Now()

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"now"}).AddRow(now)
	mock.ExpectQuery(`SELECT CURRENT_TIMESTAMP AS now`).
		WithArgs().
		WillReturnRows(rows)
	mock.ExpectCommit()

	require.NoError(t, c.Tx(context.Background(), func(tx *Tx) error {
		v, err := tx.Now(context.Background())
		require.NoError(t, err)
		require.Equal(t, now, v)

		return nil
	}))
}
