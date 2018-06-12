package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/silas/jdb"
	"github.com/silas/jdb/internal/ptr"
	"github.com/silas/jdb/test/db/internal/data"
	"github.com/stretchr/testify/require"
)

func (dt *Test) testSelect(t *testing.T) {
	dt.testSelectFirst(t)
	dt.testSelectAll(t)
	dt.testSelectCount(t)
	dt.testSelectIDs(t)
	dt.testSelectWhere(t)
	dt.testSelectOrder(t)
}

func (dt *Test) testSelectFirst(t *testing.T) {
	db := dt.setup(t, true)

	ctx := context.Background()

	require.NoError(t, db.Tx(ctx, func(tx *jdb.Tx) error {
		var user data.User
		err := db.Query(data.UserKind).Get("123").Select().First(ctx, tx, &user)
		require.Equal(t, jdb.ErrNotFound, err)

		err = db.Query(data.UserKind).Get(data.User1ID).Select().First(ctx, tx, user)
		require.EqualError(t, err, "dest must be a pointer")

		var invalidUser int
		err = db.Query(data.UserKind).Get(data.User1ID).Select().First(ctx, tx, &invalidUser)
		require.EqualError(t, err, "dest must be a struct")

		err = db.Query(data.UserKind).Get(data.User1ID).Select().First(ctx, tx, &user)
		require.NoError(t, err)
		data.RequireUser1(t, user, true)

		err = db.Query(data.UserKind).Get(data.User3ID).Select().First(ctx, tx, &user)
		require.NoError(t, err)
		data.RequireUser3(t, user)

		var userMeta data.UserMeta
		err = db.Query(data.UserKind).Get(data.User1ID).Select(db.ID, db.Data, db.CreateTime).First(ctx, tx, &userMeta)
		require.NoError(t, err)
		require.Equal(t, userMeta.ID, data.User1ID)
		require.Equal(t, userMeta.Email, data.User1Email)
		require.Equal(t, userMeta.CreateTime, data.User1CreateTime)

		var zeroTime time.Time
		err = db.Query(data.UserKind).Get(data.User1ID).Select(db.ID, db.UpdateTime).First(ctx, tx, &userMeta)
		require.NoError(t, err)
		require.Equal(t, userMeta.ID, data.User1ID)
		require.Equal(t, userMeta.Email, "")
		require.Equal(t, userMeta.CreateTime, zeroTime)

		var userPtr data.UserPtr
		err = db.Query(data.UserKind).Get(data.User1ID).Select().First(ctx, tx, &userPtr)
		require.NoError(t, err)
		data.RequireUserPtr1(t, userPtr, true)

		err = db.Query(data.UserKind).Get(data.User3ID).Select().First(ctx, tx, &userPtr)
		require.NoError(t, err)
		data.RequireUserPtr3(t, userPtr)

		var userMetaPtr data.UserMetaPtr
		err = db.Query(data.UserKind).Get(data.User1ID).Select(db.ID, db.Data, db.CreateTime).First(ctx, tx, &userMetaPtr)
		require.NoError(t, err)
		require.Equal(t, userMetaPtr.ID, ptr.String(data.User1ID))
		require.Equal(t, userMetaPtr.Email, ptr.String(data.User1Email))
		require.Equal(t, userMetaPtr.CreateTime, &data.User1CreateTime)

		err = db.Query(data.UserKind).Get(data.User1ID).Select(db.ID, db.UpdateTime).First(ctx, tx, &userMetaPtr)
		require.NoError(t, err)
		require.Equal(t, userMetaPtr.ID, ptr.String(data.User1ID))
		require.Nil(t, userMetaPtr.Email)
		require.Nil(t, userMetaPtr.CreateTime)

		return tx.Commit()
	}))
}

func (dt *Test) testSelectAll(t *testing.T) {
	db := dt.setup(t, true)

	ctx := context.Background()

	require.NoError(t, db.Tx(ctx, func(tx *jdb.Tx) error {
		var users []data.User
		err := db.Query(data.UserKind).Select().OrderBy(db.CreateTime.Asc()).All(ctx, tx, &users)
		require.NoError(t, err)
		require.Len(t, users, 3)
		data.RequireUser1(t, users[0], true)
		data.RequireUser2(t, users[1], true)
		data.RequireUser3(t, users[2])

		err = db.Query(data.UserKind).Get("-1").Select().All(ctx, tx, &users)
		require.NoError(t, err)
		require.Len(t, users, 0)

		var usersPtr []data.UserPtr
		err = db.Query(data.UserKind).Select().OrderBy(db.CreateTime.Asc()).All(ctx, tx, &usersPtr)
		require.NoError(t, err)
		require.Len(t, usersPtr, 3)
		data.RequireUserPtr1(t, usersPtr[0], true)
		data.RequireUserPtr2(t, usersPtr[1], true)
		data.RequireUserPtr3(t, usersPtr[2])

		err = db.Query(data.UserKind).Get("-1").Select().All(ctx, tx, &usersPtr)
		require.NoError(t, err)
		require.Len(t, usersPtr, 0)

		return tx.Commit()
	}))
}

func (dt *Test) testSelectCount(t *testing.T) {
	db := dt.setup(t, true)

	ctx := context.Background()

	require.NoError(t, db.Tx(ctx, func(tx *jdb.Tx) error {
		var count int
		err := db.Query(data.UserKind).Count().First(ctx, tx, &count)
		require.NoError(t, err)
		require.Equal(t, count, 3)

		err = db.Query(data.UserKind).Where(jdb.Eq(db.Path("Age"), 34)).Count().First(ctx, tx, &count)
		require.NoError(t, err)
		require.Equal(t, count, 1)

		err = db.Query(data.UserKind).Get("-1").Count().First(ctx, tx, &count)
		require.NoError(t, err)
		require.Equal(t, count, 0)

		return tx.Commit()
	}))
}

func (dt *Test) testSelectIDs(t *testing.T) {
	db := dt.setup(t, true)

	ctx := context.Background()

	require.NoError(t, db.Tx(ctx, func(tx *jdb.Tx) error {
		var ids []string
		err := db.Query(data.UserKind).Select(db.ID).OrderBy(db.CreateTime.Asc()).All(ctx, tx, &ids)
		require.NoError(t, err)
		require.Len(t, ids, 3)
		require.Equal(t, ids, []string{data.User1ID, data.User2ID, data.User3ID})

		err = db.Query(data.UserKind).Get("-1").Select(db.ID).All(ctx, tx, &ids)
		require.NoError(t, err)
		require.Len(t, ids, 0)

		return tx.Commit()
	}))
}

func (dt *Test) testSelectWhere(t *testing.T) {
	db := dt.setup(t, true)

	givenName := db.Path("Name", "GivenName")
	familyName := db.Path("Name", "FamilyName")
	age := db.Path("Age")

	tests := []struct {
		Conditions []jdb.Condition
		IDs        []string
	}{
		// Eq
		{
			[]jdb.Condition{jdb.Eq(familyName, data.User2FamilyName)},
			[]string{data.User2ID},
		},
		{
			[]jdb.Condition{jdb.Eq(familyName, "Whedon")},
			[]string{},
		},
		{
			[]jdb.Condition{jdb.Eq(familyName, nil)},
			[]string{data.User3ID},
		},
		// NotEq
		{
			[]jdb.Condition{jdb.NotEq(familyName, data.User1FamilyName)},
			[]string{data.User2ID},
		},
		{
			[]jdb.Condition{jdb.NotEq(familyName, nil)},
			[]string{data.User1ID, data.User2ID},
		},
		// Like
		{
			[]jdb.Condition{jdb.Like(givenName, "J%")},
			[]string{data.User1ID, data.User2ID},
		},
		{
			[]jdb.Condition{jdb.Like(givenName, "Z%")},
			[]string{},
		},
		{
			[]jdb.Condition{jdb.Like(givenName, nil)},
			[]string{},
		},
		// NotLike
		{
			[]jdb.Condition{jdb.NotLike(givenName, "J%")},
			[]string{},
		},
		{
			[]jdb.Condition{jdb.NotLike(familyName, "S%")},
			[]string{data.User1ID},
		},
		{
			[]jdb.Condition{jdb.NotLike(familyName, nil)},
			[]string{},
		},
		// In
		{
			[]jdb.Condition{jdb.In(age, nil, 23, 34)},
			[]string{data.User1ID, data.User2ID, data.User3ID},
		},
		{
			[]jdb.Condition{jdb.In(age, 23, 34)},
			[]string{data.User1ID, data.User2ID},
		},
		{
			[]jdb.Condition{jdb.In(age, 23)},
			[]string{data.User2ID},
		},
		{
			[]jdb.Condition{jdb.In(age, nil)},
			[]string{data.User3ID},
		},
		{
			[]jdb.Condition{jdb.In(age)},
			[]string{},
		},
		// NotIn
		{
			[]jdb.Condition{jdb.NotIn(age, nil, 23, 34)},
			[]string{},
		},
		{
			[]jdb.Condition{jdb.NotIn(age, 23, 34)},
			[]string{data.User3ID},
		},
		{
			[]jdb.Condition{jdb.NotIn(age, 23)},
			[]string{data.User1ID, data.User3ID},
		},
		{
			[]jdb.Condition{jdb.NotIn(age, nil)},
			[]string{data.User1ID, data.User2ID},
		},
		{
			[]jdb.Condition{jdb.NotIn(age)},
			[]string{data.User1ID, data.User2ID, data.User3ID},
		},
		// Gt
		{
			[]jdb.Condition{jdb.Gt(age, 30)},
			[]string{data.User1ID},
		},
		{
			[]jdb.Condition{jdb.Gt(age, 34)},
			[]string{},
		},
		{
			[]jdb.Condition{jdb.Gt(age, nil)},
			[]string{},
		},
		// Gte
		{
			[]jdb.Condition{jdb.Gte(age, 34)},
			[]string{data.User1ID},
		},
		{
			[]jdb.Condition{jdb.Gte(age, nil)},
			[]string{},
		},
		// Lt
		{
			[]jdb.Condition{jdb.Lt(age, 30)},
			[]string{data.User2ID},
		},
		{
			[]jdb.Condition{jdb.Lt(age, 23)},
			[]string{},
		},
		{
			[]jdb.Condition{jdb.Lt(age, nil)},
			[]string{},
		},
		// Lte
		{
			[]jdb.Condition{jdb.Lte(age, 23)},
			[]string{data.User2ID},
		},
		{
			[]jdb.Condition{jdb.Lte(age, nil)},
			[]string{},
		},
		// Or
		{
			[]jdb.Condition{jdb.Or(jdb.Eq(familyName, data.User1FamilyName), jdb.Eq(familyName, nil))},
			[]string{data.User1ID, data.User3ID},
		},
		{
			[]jdb.Condition{jdb.Or(jdb.Eq(familyName, "Whedon"), jdb.Gt(age, 50))},
			[]string{},
		},
		{
			[]jdb.Condition{jdb.Or()},
			[]string{data.User1ID, data.User2ID, data.User3ID},
		},
		// And
		{
			[]jdb.Condition{jdb.And(jdb.Eq(familyName, data.User1FamilyName), jdb.Eq(age, 34))},
			[]string{data.User1ID},
		},
		{
			[]jdb.Condition{jdb.And(jdb.Eq(familyName, data.User1FamilyName), jdb.Eq(familyName, data.User2FamilyName))},
			[]string{},
		},
		{
			[]jdb.Condition{jdb.And()},
			[]string{data.User1ID, data.User2ID, data.User3ID},
		},
	}

	ctx := context.Background()

	require.NoError(t, db.Tx(ctx, func(tx *jdb.Tx) error {
		for i, test := range tests {
			msg := fmt.Sprintf("Test: %d", i)

			var users []data.User
			err := db.Query(data.UserKind).
				Where(test.Conditions...).
				Select().
				OrderBy(db.CreateTime.Asc()).
				All(ctx, tx, &users)
			require.NoError(t, err, msg)

			require.Len(t, users, len(test.IDs), msg)
			for i, id := range test.IDs {
				require.Equal(t, users[i].ID, id, msg)
			}
		}

		return tx.Commit()
	}))
}

func (dt *Test) testSelectOrder(t *testing.T) {
	db := dt.setup(t, true)

	givenName := db.Path("Name", "GivenName")
	age := db.Path("Age")

	tests := []struct {
		Order []jdb.Order
		IDs   []string
	}{
		// id
		{
			[]jdb.Order{db.ID.Asc()},
			[]string{data.User1ID, data.User2ID, data.User3ID},
		},
		{
			[]jdb.Order{db.ID.Desc()},
			[]string{data.User3ID, data.User2ID, data.User1ID},
		},
		// unique string key
		{
			[]jdb.Order{db.UniqueStringKey.Asc()},
			[]string{data.User3ID, data.User1ID, data.User2ID},
		},
		{
			[]jdb.Order{db.UniqueStringKey.Desc()},
			[]string{data.User2ID, data.User1ID, data.User3ID},
		},
		// string key
		{
			[]jdb.Order{db.StringKey.Asc(), givenName.Asc()},
			[]string{data.User3ID, data.User1ID, data.User2ID},
		},
		{
			[]jdb.Order{db.StringKey.Desc(), givenName.Desc()},
			[]string{data.User2ID, data.User1ID, data.User3ID},
		},
		// numeric key
		{
			[]jdb.Order{db.NumericKey.Asc()},
			[]string{data.User3ID, data.User2ID, data.User1ID},
		},
		{
			[]jdb.Order{db.NumericKey.Desc()},
			[]string{data.User1ID, data.User2ID, data.User3ID},
		},
		// time key
		{
			[]jdb.Order{db.TimeKey.Asc()},
			[]string{data.User3ID, data.User1ID, data.User2ID},
		},
		{
			[]jdb.Order{db.TimeKey.Desc()},
			[]string{data.User2ID, data.User1ID, data.User3ID},
		},
		// path
		{
			[]jdb.Order{db.ID.Desc()},
			[]string{data.User3ID, data.User2ID, data.User1ID},
		},
		{
			[]jdb.Order{age.Desc()},
			[]string{data.User1ID, data.User2ID, data.User3ID},
		},
		// create time
		{
			[]jdb.Order{db.CreateTime.Asc()},
			[]string{data.User1ID, data.User2ID, data.User3ID},
		},
		{
			[]jdb.Order{db.CreateTime.Desc()},
			[]string{data.User3ID, data.User2ID, data.User1ID},
		},
		// update time
		{
			[]jdb.Order{db.UpdateTime.Asc()},
			[]string{data.User1ID, data.User3ID, data.User2ID},
		},
		{
			[]jdb.Order{db.UpdateTime.Desc()},
			[]string{data.User2ID, data.User3ID, data.User1ID},
		},
	}

	ctx := context.Background()

	require.NoError(t, db.Tx(ctx, func(tx *jdb.Tx) error {
		for i, test := range tests {
			msg := fmt.Sprintf("Test: %d", i)

			var users []data.User
			err := db.Query(data.UserKind).
				Select().
				OrderBy(test.Order...).
				All(ctx, tx, &users)
			require.NoError(t, err, msg)

			require.Len(t, users, len(test.IDs), msg)
			for i, id := range test.IDs {
				require.Equal(t, id, users[i].ID, msg)
			}
		}

		return tx.Commit()
	}))
}
