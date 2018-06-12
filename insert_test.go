package jdb

import (
	"testing"

	"time"

	"context"

	"github.com/silas/jdb/internal/ptr"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestInsertBuilder_ToSQL(t *testing.T) {
	c, mock := createMockClient(t)
	defer c.Close()

	kind := "kind"

	minObj := struct {
		ID string `jdb:"-id"`
	}{
		"1",
	}
	q, p, err := c.Query(kind).Insert(minObj).ToSQL()
	require.NoError(t, err)
	require.Equal(t, "INSERT INTO jdb (kind, id, parent_kind, parent_id, unique_string_key, string_key, "+
		"numeric_key, time_key, data) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", q)
	r := row{}
	require.Equal(t, params(kind, "1", r.ParentKind, r.ParentID, r.UniqueStringKey, r.StringKey, r.NumericKey,
		r.TimeKey, r.Data), p)

	tk := time.Date(2005, 3, 7, 8, 23, 34, 0, time.UTC)

	obj := struct {
		ID              string    `jdb:"-id"`
		ParentKind      string    `jdb:"-parentkind"`
		ParentId        string    `jdb:"-parentid"`
		UniqueStringKey string    `jdb:"-,uniquestringkey"`
		StringKey       string    `jdb:"-,stringkey"`
		NumericKey      float64   `jdb:"-,numerickey"`
		TimeKey         time.Time `jdb:"-,timekey"`
		Hello           string
	}{
		"2",
		"bigtest",
		"5",
		"uniqueStringKey",
		"stringKey",
		8,
		tk,
		"World",
	}
	q, p, err = c.Query(kind).Insert(obj).ToSQL()
	require.NoError(t, err)
	require.Equal(t, "INSERT INTO jdb (kind, id, parent_kind, parent_id, unique_string_key, string_key, "+
		"numeric_key, time_key, data) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", q)
	require.Equal(t, params(kind, "2", ptr.String("bigtest"), ptr.String("5"), ptr.String("uniqueStringKey"),
		ptr.String("stringKey"), ptr.Float64(8), &tk, ptr.String(`{"Hello":"World"}`)), p)

	q, p, err = c.Query(kind).Insert(minObj, obj).ToSQL()
	require.NoError(t, err)
	require.Equal(t, "INSERT INTO jdb (kind, id, parent_kind, parent_id, unique_string_key, string_key, "+
		"numeric_key, time_key, data) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?, ?)", q)
	require.Equal(t, params(kind, "1", r.ParentKind, r.ParentID, r.UniqueStringKey, r.StringKey, r.NumericKey,
		r.TimeKey, r.Data, kind, "2", ptr.String("bigtest"), ptr.String("5"), ptr.String("uniqueStringKey"),
		ptr.String("stringKey"), ptr.Float64(8), &tk, ptr.String(`{"Hello":"World"}`)), p)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertBuilder_Exec(t *testing.T) {
	c, mock := createMockClient(t)
	defer c.Close()

	kind := "test"

	tk := time.Date(2005, 3, 7, 8, 23, 34, 0, time.UTC)
	obj := struct {
		ID              string    `jdb:"-id"`
		ParentKind      string    `jdb:"-parentkind"`
		ParentId        string    `jdb:"-parentid"`
		UniqueStringKey string    `jdb:"-,uniquestringkey"`
		StringKey       string    `jdb:"-,stringkey"`
		NumericKey      float64   `jdb:"-,numerickey"`
		TimeKey         time.Time `jdb:"-,timekey"`
		Hello           string
	}{
		"2",
		"bigtest",
		"5",
		"uniqueStringKey",
		"stringKey",
		8,
		tk,
		"World",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO jdb").
		WithArgs(kind, "2", ptr.String("bigtest"), ptr.String("5"), ptr.String("uniqueStringKey"),
			ptr.String("stringKey"), ptr.Float64(8), &tk, ptr.String(`{"Hello":"World"}`)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	c.Tx(context.Background(), func(tx *Tx) error {
		err := c.Query(kind).Insert(obj).Exec(context.Background(), tx)
		require.NoError(t, err)

		return tx.Commit()
	})

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertBuilder_Add(t *testing.T) {
	q := setupQuery(t)

	i := newInsertBuilder(q)
	require.Nil(t, i.values)

	i2 := i.Add()
	require.Nil(t, i.values)
	require.Nil(t, i2.values)

	i3 := i2.Add(1)
	require.Nil(t, i.values)
	require.Nil(t, i2.values)
	require.Equal(t, i3.values, []interface{}{1})

	i4 := i3.Add(2, 3)
	require.Nil(t, i.values)
	require.Nil(t, i2.values)
	require.Equal(t, i3.values, []interface{}{1})
	require.Equal(t, i4.values, []interface{}{1, 2, 3})
}
