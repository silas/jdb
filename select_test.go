package jdb

import (
	"testing"

	"fmt"

	"context"

	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestSelectBuilder_ToSQL(t *testing.T) {
	c, mock := createMockClient(t)
	defer c.Close()

	kind := "test"
	columns := "kind, id, parent_kind, parent_id, data, create_time, update_time"
	from := "FROM jdb"
	where := "WHERE ((kind = ?))"
	query := fmt.Sprintf("SELECT %s %s %s", columns, from, where)

	tests := []struct {
		Builder QueryBuilder
		Query   string
		Params  []interface{}
	}{
		{
			c.Query(kind).Select(),
			query,
			params(kind),
		},
		{
			c.Query(kind).Select(c.ID),
			fmt.Sprintf("SELECT id %s %s", from, where),
			params(kind),
		},
		{
			c.Query(kind).Select(c.ID, c.CreateTime),
			fmt.Sprintf("SELECT id, create_time %s %s", from, where),
			params(kind),
		},
		{
			c.Query(kind).Select().Limit(10),
			query + " LIMIT 10",
			params(kind),
		},
		{
			c.Query(kind).Select().Offset(5),
			query + " OFFSET 5",
			params(kind),
		},
		{
			c.Query(kind).Select().OrderBy(c.NumericKey.Asc()),
			query + " ORDER BY numeric_key ASC",
			params(kind),
		},
		{
			c.Query(kind).Select().OrderBy(c.NumericKey.Desc(), c.ID.Asc()),
			query + " ORDER BY numeric_key DESC, id ASC",
			params(kind),
		},
		{
			c.Query(kind).Where(Eq(c.StringKey, "example.com")).Select(),
			fmt.Sprintf("SELECT %s %s WHERE ((kind = ?) AND (string_key = ?))", columns, from),
			params(kind, "example.com"),
		},
	}

	for i, test := range tests {
		msg := fmt.Sprintf("Test: %d", i)

		s, p, err := test.Builder.ToSQL()
		require.NoError(t, err, msg)
		require.Equal(t, test.Query, s, msg)
		require.Equal(t, test.Params, p, msg)
	}

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSelectBuilder_Exec(t *testing.T) {
	c, mock := createMockClient(t)
	defer c.Close()

	createTime := time.Date(2005, 3, 7, 8, 23, 34, 0, time.UTC)
	updateTime := time.Date(2018, 5, 22, 1, 5, 2, 0, time.UTC)

	type obj struct {
		Kind       string    `jdb:"-kind"`
		ID         string    `jdb:"-id"`
		ParentKind string    `jdb:"-parentkind"`
		ParentID   string    `jdb:"-parentid"`
		CreateTime time.Time `jdb:"-createtime"`
		UpdateTime time.Time `jdb:"-updatetime"`
		Hello      string
	}

	kind := "test"
	columns := []string{"kind", "id", "parent_kind", "parent_id", "data", "create_time", "update_time"}

	mock.ExpectBegin()
	firstRows := sqlmock.NewRows(columns).
		AddRow(kind, "1", "parentKind", "parentID", `{"Hello":"World"}`, createTime, updateTime)
	mock.ExpectQuery(`SELECT kind, id, .*id =.*`).
		WithArgs(kind, "1").
		WillReturnRows(firstRows)
	allRows := sqlmock.NewRows(columns).
		AddRow(kind, "1", "parentKind", "parentID", `{"Hello":"World"}`, createTime, updateTime).
		AddRow(kind, "2", nil, nil, `{"Hello":"World 2"}`, createTime.AddDate(1, 0, 0),
			updateTime.AddDate(1, 0, 0))
	mock.ExpectQuery(`SELECT kind, id, .*string_key =.*`).
		WithArgs(kind, "test").
		WillReturnRows(allRows)
	mock.ExpectCommit()

	c.Tx(context.Background(), func(tx *Tx) error {
		var result obj
		err := c.Query(kind).Get("1").Select().First(context.Background(), tx, &result)
		require.NoError(t, err)

		require.Equal(t, kind, result.Kind)
		require.Equal(t, "1", result.ID)
		require.Equal(t, "parentKind", result.ParentKind)
		require.Equal(t, "parentID", result.ParentID)
		require.Equal(t, createTime, result.CreateTime)
		require.Equal(t, updateTime, result.UpdateTime)
		require.Equal(t, "World", result.Hello)

		var results []obj
		err = c.Query(kind).Where(Eq(c.StringKey, "test")).Select().All(context.Background(), tx, &results)
		require.NoError(t, err)
		require.Len(t, results, 2)
		// result 1
		require.Equal(t, kind, results[0].Kind)
		require.Equal(t, "1", results[0].ID)
		require.Equal(t, "parentKind", results[0].ParentKind)
		require.Equal(t, "parentID", results[0].ParentID)
		require.Equal(t, createTime, results[0].CreateTime)
		require.Equal(t, updateTime, results[0].UpdateTime)
		require.Equal(t, "World", results[0].Hello)
		// result 2
		require.Equal(t, kind, results[1].Kind)
		require.Equal(t, "2", results[1].ID)
		require.Equal(t, "", results[1].ParentKind)
		require.Equal(t, "", results[1].ParentID)
		require.Equal(t, createTime.AddDate(1, 0, 0), results[1].CreateTime)
		require.Equal(t, updateTime.AddDate(1, 0, 0), results[1].UpdateTime)
		require.Equal(t, "World 2", results[1].Hello)

		return tx.Commit()
	})

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSelectBuilder_Limit(t *testing.T) {
	s := setupQuery(t).Select()
	require.False(t, s.limitDefined)
	require.Zero(t, s.limit)

	s2 := s.Limit(10)
	require.False(t, s.limitDefined)
	require.Zero(t, s.limit)
	require.True(t, s2.limitDefined)
	require.Equal(t, uint64(10), s2.limit)

	s3 := s2.Limit(20)
	require.False(t, s.limitDefined)
	require.Zero(t, s.limit)
	require.True(t, s2.limitDefined)
	require.Equal(t, uint64(10), s2.limit)
	require.True(t, s3.limitDefined)
	require.Equal(t, uint64(20), s3.limit)
}

func TestSelectBuilder_Offset(t *testing.T) {
	s := setupQuery(t).Select()
	require.False(t, s.offsetDefined)
	require.Zero(t, s.offset)

	s2 := s.Offset(10)
	require.False(t, s.offsetDefined)
	require.Zero(t, s.offset)
	require.True(t, s2.offsetDefined)
	require.Equal(t, uint64(10), s2.offset)

	s3 := s2.Offset(20)
	require.False(t, s.offsetDefined)
	require.Zero(t, s.offset)
	require.True(t, s2.offsetDefined)
	require.Equal(t, uint64(10), s2.offset)
	require.True(t, s3.offsetDefined)
	require.Equal(t, uint64(20), s3.offset)
}

func TestSelectBuilder_OrderBy(t *testing.T) {
	s := setupQuery(t).Select()
	require.Nil(t, s.order)

	s2 := s.OrderBy(createTimeField.Asc())
	require.Nil(t, s.order)
	require.Equal(t, []Order{{createTimeField, false}}, s2.order)

	s3 := s2.OrderBy(stringKeyField.Desc())
	require.Nil(t, s.order)
	require.Equal(t, []Order{{createTimeField, false}}, s2.order)
	require.Equal(t, []Order{{createTimeField, false}, {stringKeyField, true}}, s3.order)
}
