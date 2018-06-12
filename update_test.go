package jdb

import (
	"testing"
	"time"

	"context"

	"github.com/silas/jdb/internal/ptr"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestUpdateBuilder_ToSQL(t *testing.T) {
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

	s, p, err := c.Query(kind).Update(obj).ToSQL()
	require.NoError(t, err)
	require.Equal(t, "UPDATE jdb SET parent_kind = ?, parent_id = ?, unique_string_key = ?, string_key = ?, "+
		"numeric_key = ?, time_key = ?, data = ?, update_time = CURRENT_TIMESTAMP WHERE ((kind = ?) AND (id = ?))", s)
	require.Equal(t, params(ptr.String(obj.ParentKind), ptr.String(obj.ParentId), ptr.String(obj.UniqueStringKey),
		ptr.String(obj.StringKey), ptr.Float64(obj.NumericKey), &obj.TimeKey, ptr.String(`{"Hello":"World"}`),
		kind, obj.ID), p)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateBuilder_Exec(t *testing.T) {
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
	mock.ExpectExec("UPDATE jdb SET").
		WithArgs(ptr.String(obj.ParentKind), ptr.String(obj.ParentId), ptr.String(obj.UniqueStringKey),
			ptr.String(obj.StringKey), ptr.Float64(obj.NumericKey), &obj.TimeKey, ptr.String(`{"Hello":"World"}`),
			kind, obj.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	c.Tx(context.Background(), func(tx *Tx) error {
		err := c.Query(kind).Update(obj).Exec(context.Background(), tx)
		require.NoError(t, err)

		return tx.Commit()
	})

	require.NoError(t, mock.ExpectationsWereMet())
}
