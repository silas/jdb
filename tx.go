package jdb

import (
	"context"
	"database/sql"
	"time"
)

type Tx struct {
	c  *Client
	tx *sql.Tx
}

func (t *Tx) Commit() error {
	return t.tx.Commit()
}

func (t *Tx) exec(ctx context.Context, builder QueryBuilder) (sql.Result, error) {
	if t.c.readOnly {
		return nil, ErrReadOnlyMode
	}

	query, params, err := builder.ToSQL()
	if err != nil {
		return nil, err
	}

	result, err := t.tx.ExecContext(ctx, query, params...)
	if err != nil {
		err = t.c.d.ErrorMap(err)
	}
	return result, err
}

func (t *Tx) query(ctx context.Context, builder *SelectBuilder) (*Rows, error) {
	query, params, err := builder.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := t.tx.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, t.c.d.ErrorMap(err)
	}
	return newRows(rows, builder.columns), nil
}

func (t *Tx) Now(ctx context.Context) (time.Time, error) {
	now, err := t.c.d.Now(ctx, t.tx)
	if err != nil {
		err = t.c.d.ErrorMap(err)
	}
	return now, err
}
