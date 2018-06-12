package jdb

import (
	"bytes"
	"context"
)

var updateColumnsSQL = "parent_kind = ?, parent_id = ?, unique_string_key = ?, string_key = ?, numeric_key = ?, time_key = ?, data = ?, update_time = "

type UpdateBuilder struct {
	q  *Query
	wb *WhereBuilder

	value interface{}
}

func newUpdateBuilder(q *Query, wb *WhereBuilder, value interface{}) *UpdateBuilder {
	if wb == nil {
		wb = newWhereBuilder(q)
	}
	return &UpdateBuilder{q: q, wb: wb, value: value}
}

func (b *UpdateBuilder) Exec(ctx context.Context, tx *Tx) error {
	_, err := tx.exec(ctx, b)
	return err
}

func (b *UpdateBuilder) ToSQL() (string, []interface{}, error) {
	var params []interface{}
	query := &bytes.Buffer{}

	r, err := rowScanInput(b.q.kind, b.value)
	if err != nil {
		return "", nil, err
	}

	query.WriteString("UPDATE ")
	query.WriteString(b.q.table)
	query.WriteString(" SET ")

	query.WriteString(updateColumnsSQL)
	query.WriteString(b.q.d.TimestampExpression())
	query.WriteString(" ")
	params = append(params, r.ParentKind, r.ParentID, r.UniqueStringKey, r.StringKey, r.NumericKey, r.TimeKey, r.Data)

	err = b.wb.toWhereSQL(query, &params)
	if err != nil {
		return "", nil, err
	}

	return b.q.d.ReplacePlaceHolders(query.String()), params, nil
}
