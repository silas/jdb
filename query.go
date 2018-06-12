package jdb

import (
	"github.com/silas/jdb/dialect"
)

type QueryBuilder interface {
	ToSQL() (string, []interface{}, error)
}

type Query struct {
	d     dialect.Dialect
	table string
	kind  string
}

func newQuery(d dialect.Dialect, table, kind string) *Query {
	return &Query{
		d:     d,
		table: table,
		kind:  kind,
	}
}

func (q *Query) get(ids ...string) *WhereBuilder {
	if len(ids) == 0 {
		return q.Where()
	}
	if len(ids) == 1 {
		return q.Where(Eq(idField, ids[0]))
	}
	tmp := make([]interface{}, len(ids))
	for i := range ids {
		tmp[i] = ids[i]
	}
	return q.Where(In(idField, tmp...))
}

func (q *Query) Get(ids ...string) *WhereBuilder {
	return q.get(ids...)
}

func (q *Query) Select(columns ...SelectField) *SelectBuilder {
	if len(columns) == 0 {
		columns = defaultSelectColumns
	}
	return newSelectBuilder(q, nil, columns)
}

func (q *Query) Count() *SelectBuilder {
	return newSelectBuilder(q, nil, []SelectField{selectCount{}})
}

func (q *Query) Where(where ...Condition) *WhereBuilder {
	return newWhereBuilder(q).Where(where...)
}

func (q *Query) Delete(ids ...string) *DeleteBuilder {
	return q.get(ids...).Delete()
}

func (q *Query) Insert(values ...interface{}) *InsertBuilder {
	return newInsertBuilder(q).Add(values...)
}

func (q *Query) Update(value interface{}) *UpdateBuilder {
	kind, id, err := idScanInput(value)
	if err != nil || (kind != "" && kind != q.kind) {
		return q.Where(False()).update(value)
	} else {
		return q.get(id).update(value)
	}
}
