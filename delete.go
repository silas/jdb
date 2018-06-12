package jdb

import (
	"bytes"
	"context"
)

type DeleteBuilder struct {
	q *Query

	wb *WhereBuilder
}

func newDeleteBuilder(q *Query, wb *WhereBuilder) *DeleteBuilder {
	if wb == nil {
		wb = newWhereBuilder(q)
	}
	return &DeleteBuilder{q: q, wb: wb}
}

func (b *DeleteBuilder) ToSQL() (string, []interface{}, error) {
	var params []interface{}
	query := &bytes.Buffer{}

	query.WriteString("DELETE FROM ")
	query.WriteString(b.q.table)
	query.WriteString(" ")

	err := b.wb.toWhereSQL(query, &params)
	if err != nil {
		return "", nil, err
	}

	return b.q.d.ReplacePlaceHolders(query.String()), params, nil
}

func (b *DeleteBuilder) Exec(ctx context.Context, tx *Tx) error {
	_, err := tx.exec(ctx, b)
	return err
}
