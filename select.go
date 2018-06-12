package jdb

import (
	"bytes"
	"context"
	"strconv"
)

type SelectBuilder struct {
	q  *Query
	wb *WhereBuilder

	columns       []SelectField
	limit         uint64
	limitDefined  bool
	offset        uint64
	offsetDefined bool
	order         []Order
}

type selectCount struct{}

var defaultSelectColumns = []SelectField{
	kindField, idField, parentKindField, parentIdField, dataField, createTimeField, updateTimeField}

func (sc selectCount) toSelectField() string {
	return "count(*) AS count"
}

func newSelectBuilder(q *Query, wb *WhereBuilder, columns []SelectField) *SelectBuilder {
	if wb == nil {
		wb = newWhereBuilder(q)
	}
	return &SelectBuilder{q: q, wb: wb, columns: columns}
}

func (b *SelectBuilder) Limit(v uint64) *SelectBuilder {
	if b.limit == v {
		return b
	}
	n := *b
	n.limit = v
	n.limitDefined = true
	return &n
}

func (b *SelectBuilder) Offset(v uint64) *SelectBuilder {
	if b.offset == v {
		return b
	}
	n := *b
	n.offset = v
	n.offsetDefined = true
	return &n
}

func (b *SelectBuilder) OrderBy(order ...Order) *SelectBuilder {
	if len(order) == 0 {
		return b
	}
	n := *b
	if len(b.order) == 0 {
		n.order = order
		return &n
	}

	total := len(b.order) + len(order)
	tmp := make([]Order, total)
	start := copy(tmp, b.order)
	for i, v := range order {
		tmp[i+start] = v
	}
	n.order = tmp

	return &n
}

func (b *SelectBuilder) ToSQL() (string, []interface{}, error) {
	var params []interface{}
	query := &bytes.Buffer{}

	query.WriteString("SELECT ")
	for i, c := range b.columns {
		if i != 0 {
			query.WriteString(", ")
		}
		query.WriteString(c.toSelectField())
	}
	query.WriteString(" ")
	query.WriteString("FROM ")
	query.WriteString(b.q.table)
	query.WriteString(" ")

	err := b.wb.toWhereSQL(query, &params)
	if err != nil {
		return "", nil, err
	}

	for i, order := range b.order {
		if i == 0 {
			query.WriteString(" ORDER BY ")
		} else {
			query.WriteString(", ")
		}
		query.WriteString(b.q.d.OrderExpression(order))
	}

	if b.limitDefined {
		query.WriteString(" LIMIT ")
		query.WriteString(strconv.FormatUint(b.limit, 10))
	}

	if b.offsetDefined {
		query.WriteString(" OFFSET ")
		query.WriteString(strconv.FormatUint(b.offset, 10))
	}

	return b.q.d.ReplacePlaceHolders(query.String()), params, nil
}

func (b *SelectBuilder) Rows(ctx context.Context, tx *Tx) (*Rows, error) {
	return tx.query(ctx, b)
}

func (b *SelectBuilder) First(ctx context.Context, tx *Tx, dest interface{}) error {
	rows, err := b.Rows(ctx, tx)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}
		return ErrNotFound
	}

	return rows.Scan(dest)
}

func (b *SelectBuilder) All(ctx context.Context, tx *Tx, dest interface{}) error {
	rows, err := b.Rows(ctx, tx)
	if err != nil {
		return err
	}
	defer rows.Close()

	return rows.ScanAll(dest)
}
