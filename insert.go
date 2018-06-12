package jdb

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
)

const insertColumnsSQL = "kind, id, parent_kind, parent_id, unique_string_key, string_key, numeric_key, time_key, data"
const insertPlaceholdersSQL = "?, ?, ?, ?, ?, ?, ?, ?, ?"

type InsertBuilder struct {
	q *Query

	values []interface{}
}

func newInsertBuilder(q *Query) *InsertBuilder {
	return &InsertBuilder{q: q}
}

func (b *InsertBuilder) Add(values ...interface{}) *InsertBuilder {
	if len(values) == 0 {
		return b
	}
	n := *b
	if len(b.values) == 0 {
		n.values = values
		return &n
	}

	total := len(b.values) + len(values)
	tmp := make([]interface{}, total)
	start := copy(tmp, b.values)
	for i, v := range values {
		tmp[i+start] = v
	}
	n.values = tmp

	return &n
}

func (b *InsertBuilder) Exec(ctx context.Context, tx *Tx) error {
	_, err := tx.exec(ctx, b)
	return err
}

func (b *InsertBuilder) ToSQL() (string, []interface{}, error) {
	var params []interface{}
	query := &bytes.Buffer{}

	query.WriteString("INSERT INTO ")
	query.WriteString(b.q.table)
	query.WriteString(" (")
	query.WriteString(insertColumnsSQL)
	query.WriteString(") VALUES")

	for i, v := range b.values {
		r, err := rowScanInput(b.q.kind, v)
		if err != nil {
			return "", nil, fmt.Errorf("value %d %s", i, err)
		}

		if i > 0 {
			query.WriteString(",")
		}
		query.WriteString(" (")
		query.WriteString(insertPlaceholdersSQL)
		query.WriteString(")")
		params = append(params, r.Kind, r.ID, r.ParentKind, r.ParentID, r.UniqueStringKey, r.StringKey, r.NumericKey,
			r.TimeKey, r.Data)
	}

	return b.q.d.ReplacePlaceHolders(query.String()), params, nil
}

func isZero(t reflect.Type, v reflect.Value) bool {
	zero := reflect.Zero(t).Interface()
	return reflect.DeepEqual(v, zero)
}
