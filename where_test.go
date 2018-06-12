package jdb

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWhere(t *testing.T) {
	c, mock := createMockClient(t)
	defer c.Close()

	var p []interface{}
	q := &bytes.Buffer{}

	err := c.Query("test").Where(Eq(idField, "1")).toWhereSQL(q, &p)
	require.NoError(t, err)
	require.Equal(t, "WHERE ((kind = ?) AND (id = ?))", q.String())
	require.Equal(t, params("test", "1"), p)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestWhereBuilder_Where(t *testing.T) {
	q := setupQuery(t)

	w := newWhereBuilder(q)
	require.Equal(t, And(Eq(kindField, w.q.kind)), w.where)

	w2 := w.Where(True())
	require.Equal(t, And(Eq(kindField, w.q.kind)), w.where)
	require.Equal(t, And(Eq(kindField, w.q.kind), True()), w2.where)

	w3 := w2.Where(False())
	require.Equal(t, And(Eq(kindField, w.q.kind)), w.where)
	require.Equal(t, And(Eq(kindField, w.q.kind), True()), w2.where)
	require.Equal(t, And(Eq(kindField, w.q.kind), True(), False()), w3.where)
}
