package data

import (
	"testing"

	"github.com/silas/jdb/internal/ptr"
	"github.com/stretchr/testify/require"
)

func RequireUser1(t *testing.T, user User, checkTime bool) {
	require.Equal(t, user.ID, User1ID)
	require.Equal(t, user.Email, User1Email)
	require.Equal(t, user.Name.GivenName, User1GivenName)
	require.Equal(t, user.Name.FamilyName, User1FamilyName)
	require.Equal(t, user.Name.Aliases, User1Aliases)
	require.Equal(t, user.Age, User1Age)
	if checkTime {
		require.Equal(t, user.CreateTime, User1CreateTime)
		require.Equal(t, user.UpdateTime, User1UpdateTime)
	}
}

func RequireUserPtr1(t *testing.T, user UserPtr, checkTime bool) {
	var aliases []*string
	for _, v := range User1Aliases {
		aliases = append(aliases, ptr.String(v))
	}

	require.Equal(t, user.ID, ptr.String(User1ID))
	require.Equal(t, user.Email, ptr.String(User1Email))
	require.Equal(t, user.Name.GivenName, ptr.String(User1GivenName))
	require.Equal(t, user.Name.FamilyName, ptr.String(User1FamilyName))
	require.Equal(t, user.Name.Aliases, &aliases)
	require.Equal(t, user.Age, ptr.Int(User1Age))
	if checkTime {
		require.Equal(t, user.CreateTime, &User1CreateTime)
		require.Equal(t, user.UpdateTime, &User1UpdateTime)
	}
}

func RequireUser2(t *testing.T, user User, checkTime bool) {
	require.Equal(t, user.ID, User2ID)
	require.Equal(t, user.Email, User2Email)
	require.Equal(t, user.Name.GivenName, User2GivenName)
	require.Equal(t, user.Name.FamilyName, User2FamilyName)
	require.Equal(t, user.Name.Aliases, User2Aliases)
	require.Equal(t, user.Age, User2Age)
	if checkTime {
		require.Equal(t, user.CreateTime, User2CreateTime)
		require.Equal(t, user.UpdateTime, User2UpdateTime)
	}
}

func RequireUserPtr2(t *testing.T, user UserPtr, checkTime bool) {
	var aliases []*string
	for _, v := range User2Aliases {
		aliases = append(aliases, ptr.String(v))
	}

	require.Equal(t, user.ID, ptr.String(User2ID))
	require.Equal(t, user.Email, ptr.String(User2Email))
	require.Equal(t, user.Name.GivenName, ptr.String(User2GivenName))
	require.Equal(t, user.Name.FamilyName, ptr.String(User2FamilyName))
	require.Equal(t, user.Name.Aliases, &aliases)
	require.Equal(t, user.Age, ptr.Int(User2Age))
	if checkTime {
		require.Equal(t, user.CreateTime, &User2CreateTime)
		require.Equal(t, user.UpdateTime, &User2UpdateTime)
	}
}

func RequireUser3(t *testing.T, user User) {
	require.Equal(t, user.ID, User3ID)
	require.Equal(t, user.Email, "")
	require.Equal(t, user.Name.GivenName, "")
	require.Equal(t, user.Name.FamilyName, "")
	require.Equal(t, user.Age, 0)
}

func RequireUserPtr3(t *testing.T, user UserPtr) {
	require.Equal(t, user.ID, ptr.String(User3ID))
	require.Nil(t, user.Email)
	require.Nil(t, user.Name)
	require.Nil(t, user.Age)
}
