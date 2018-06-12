package db

import (
	"context"
	"testing"

	"github.com/silas/jdb"
	"github.com/silas/jdb/internal/ptr"
	"github.com/silas/jdb/test/db/internal/data"
	"github.com/stretchr/testify/require"
)

func (dt *Test) testInsert(t *testing.T) {
	db := dt.setup(t, false)

	ctx := context.Background()
	query := db.Query("user")

	require.NoError(t, db.Tx(ctx, func(tx *jdb.Tx) error {
		inputUser := data.User{
			ID:    data.User1ID,
			Email: data.User1Email,
			Name: data.Name{
				GivenName:  data.User1GivenName,
				FamilyName: data.User1FamilyName,
				Aliases:    data.User1Aliases,
			},
			Age: data.User1Age,
		}

		err := query.Insert(inputUser).Exec(ctx, tx)
		require.NoError(t, err)

		var user data.User
		err = query.Get(data.User1ID).Select().First(ctx, tx, &user)
		require.NoError(t, err)
		data.RequireUser1(t, user, false)

		query.Delete(data.User1ID).Exec(ctx, tx)

		var aliases []*string
		for _, v := range data.User1Aliases {
			aliases = append(aliases, ptr.String(v))
		}

		inputUserPtr := &data.UserPtr{
			ID:    ptr.String(data.User1ID),
			Email: ptr.String(data.User1Email),
			Name: &data.NamePtr{
				GivenName:  ptr.String(data.User1GivenName),
				FamilyName: ptr.String(data.User1FamilyName),
				Aliases:    &aliases,
			},
			Age: ptr.Int(data.User1Age),
		}

		err = query.Insert(inputUserPtr).Exec(ctx, tx)
		require.NoError(t, err)

		var userPtr data.UserPtr
		err = query.Get(data.User1ID).Select().First(ctx, tx, &userPtr)
		require.NoError(t, err)
		data.RequireUserPtr1(t, userPtr, false)

		return tx.Commit()
	}))
}
