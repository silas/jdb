package jdb

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/silas/jdb/dialect"
)

var stringValue reflect.Value

func init() {
	stringValue = reflect.ValueOf("")
}

type Client struct {
	d        dialect.Dialect
	db       *sql.DB
	table    string
	readOnly bool

	ID              SelectWhereColumn
	Kind            SelectWhereColumn
	ParentKind      SelectWhereColumn
	ParentId        SelectWhereColumn
	UniqueStringKey WhereColumn
	StringKey       WhereColumn
	NumericKey      WhereColumn
	TimeKey         WhereColumn
	Data            SelectColumn
	CreateTime      SelectWhereColumn
	UpdateTime      SelectWhereColumn
}

func Open(driverName, dataSourceName string, opts ...Option) (*Client, error) {
	table := "jdb"
	readOnly := false

	for _, opt := range opts {
		switch v := opt.(type) {
		case optionTable:
			if v.table == "" {
				return nil, fmt.Errorf("jdb: invalid table name")
			}
			table = v.table
		case optionReadOnly:
			readOnly = v.readOnly
		default:
			panic("unknown option")
		}
	}

	d, err := Dialect(driverName)
	if err != nil {
		return nil, err
	}

	dataSourceNameOpts := dialect.ValidateDataSourceNameOpts{
		ReadOnly: readOnly,
	}
	err = d.ValidateDataSourceName(dataSourceName, dataSourceNameOpts)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	c := &Client{
		d:        d,
		db:       db,
		table:    table,
		readOnly: readOnly,

		ID:              idField,
		ParentKind:      parentKindField,
		ParentId:        parentIdField,
		UniqueStringKey: uniqueStringKeyField,
		StringKey:       stringKeyField,
		NumericKey:      numericKeyField,
		TimeKey:         timeKeyField,
		Data:            dataField,
		CreateTime:      createTimeField,
		UpdateTime:      updateTimeField,
	}

	return c, nil
}

func (c *Client) Migrate(ctx context.Context) error {
	return c.d.Migrate(ctx, c.db, c.table)
}

func (c *Client) Path(key ...string) *PathField {
	p := &PathField{c.d.Path()}
	if len(key) > 0 {
		for _, v := range key {
			p.Key(v)
		}
	}
	return p
}

func (c *Client) Query(kind string) *Query {
	return newQuery(c.d, c.table, kind)
}

func (c *Client) Tx(ctx context.Context, fn func(*Tx) error) error {
	opts := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  c.readOnly,
	}
	tx, err := c.db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	t := &Tx{c: c, tx: tx}

	return fn(t)
}

func (c *Client) Close() error {
	return c.db.Close()
}
