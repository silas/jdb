package db

import (
	"context"
	"strings"
	"testing"

	"github.com/silas/jdb"
	"github.com/silas/jdb/test/db/internal/data"
	"github.com/stretchr/testify/require"
)

func (dt *Test) testUpdate(t *testing.T) {
	db := dt.setup(t, true)

	ctx := context.Background()
	query := db.Query("user")

	require.NoError(t, db.Tx(ctx, func(tx *jdb.Tx) error {
		var user data.User
		err := query.Get(data.User1ID).Select().First(ctx, tx, &user)
		require.NoError(t, err)
		data.RequireUser1(t, user, true)

		domain := "example.org"
		email := "jane@" + strings.ToUpper(domain)
		aliases := append(user.Name.Aliases, "JR")

		user.Email = email
		user.Name.Aliases = aliases

		err = query.Update(user).Exec(ctx, tx)
		require.NoError(t, err)

		var freshUser data.User
		err = query.Get(data.User1ID).Select().First(ctx, tx, &freshUser)
		require.NoError(t, err)
		require.Equal(t, email, freshUser.Email)
		require.Equal(t, aliases, freshUser.Name.Aliases)
		require.Equal(t, data.User1GivenName, freshUser.Name.GivenName)
		require.Equal(t, data.User1FamilyName, freshUser.Name.FamilyName)
		require.Equal(t, data.User1Age, freshUser.Age)
		require.Equal(t, data.User1CreateTime, freshUser.CreateTime)
		require.True(t, freshUser.UpdateTime.After(data.User1UpdateTime))

		var users []data.User
		err = query.Select().All(ctx, tx, &users)
		require.NoError(t, err)
		require.Len(t, users, 3)
		for _, u := range users {
			if u.ID == data.User1ID {
				require.Equal(t, email, u.Email)
			} else {
				require.NotEqual(t, email, u.Email)
			}
		}

		return tx.Commit()
	}))
}
