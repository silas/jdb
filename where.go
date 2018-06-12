package jdb

import (
	"bytes"
)

type WhereBuilder struct {
	q *Query

	where and
}

func newWhereBuilder(q *Query) *WhereBuilder {
	b := &WhereBuilder{q: q}
	return b.Where(Eq(kindField, q.kind))
}

func (b *WhereBuilder) Where(where ...Condition) *WhereBuilder {
	if len(where) == 0 {
		return b
	}
	n := *b
	if len(b.where) == 0 {
		n.where = and(where)
		return &n
	}

	total := len(b.where) + len(where)
	tmp := make(and, total)
	start := copy(tmp, b.where)
	for i, v := range where {
		tmp[i+start] = v
	}
	n.where = tmp

	return &n
}

func (b *WhereBuilder) toWhereSQL(query *bytes.Buffer, params *[]interface{}) error {
	query.WriteString("WHERE ")
	return b.where.toConditionSQL(query, params)
}

func (b *WhereBuilder) Delete() *DeleteBuilder {
	return newDeleteBuilder(b.q, b)
}

func (b *WhereBuilder) Select(columns ...SelectField) *SelectBuilder {
	if len(columns) == 0 {
		columns = defaultSelectColumns
	}
	return newSelectBuilder(b.q, b, columns)
}

func (b *WhereBuilder) Count() *SelectBuilder {
	return newSelectBuilder(b.q, b, []SelectField{selectCount{}})
}

func (b *WhereBuilder) update(value interface{}) *UpdateBuilder {
	return newUpdateBuilder(b.q, b, value)
}
