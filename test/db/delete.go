package db

import (
	"context"
	"testing"

	"github.com/silas/jdb"
	"github.com/silas/jdb/test/db/internal/data"
	"github.com/stretchr/testify/require"
)

func (dt *Test) testDelete(t *testing.T) {
	db := dt.setup(t, true)

	ctx := context.Background()

	requireCount := func(q *jdb.DeleteBuilder, count int) {
		require.NoError(t, db.Tx(ctx, func(tx *jdb.Tx) error {
			q.Exec(ctx, tx)

			return tx.Commit()
		}))

		require.Equal(t, count, dt.count(t, data.UserKind))
	}

	requireCount(db.Query("nope").Where(jdb.Eq(db.ID, data.User1ID)).Delete(), 3)
	requireCount(db.Query(data.UserKind).Where(jdb.Eq(db.ID, "123")).Delete(), 3)
	requireCount(db.Query(data.UserKind).Where(jdb.Eq(db.ID, data.User1ID)).Delete(), 2)
}
