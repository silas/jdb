package db

import (
	"context"
	"testing"
	"time"

	"github.com/silas/jdb"
	"github.com/stretchr/testify/require"
)

func (dt *Test) testClient(t *testing.T) {
	db := dt.setup(t, false)

	ctx := context.Background()

	require.NoError(t, db.Tx(ctx, func(tx *jdb.Tx) error {
		now, err := tx.Now(ctx)
		require.NoError(t, err)
		require.WithinDuration(t, time.Now(), now, time.Minute)

		return tx.Commit()
	}))
}
