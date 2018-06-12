package jdb

import "github.com/silas/jdb/dialect"

type SelectField interface {
	toSelectField() string
}

type WhereField interface {
	toWhereField() string
	Asc() Order
	Desc() Order
}

type SelectColumn struct {
	n string
}

func (c SelectColumn) toSelectField() string {
	return c.n
}

type SelectWhereColumn struct {
	n string
}

func (c SelectWhereColumn) toSelectField() string {
	return c.n
}

func (c SelectWhereColumn) toWhereField() string {
	return c.n
}

func (c SelectWhereColumn) Asc() Order {
	return Order{c, false}
}

func (c SelectWhereColumn) Desc() Order {
	return Order{c, true}
}

type WhereColumn struct {
	n string
}

func (c WhereColumn) toWhereField() string {
	return c.n
}

func (c WhereColumn) Asc() Order {
	return Order{c, false}
}

func (c WhereColumn) Desc() Order {
	return Order{c, true}
}

var (
	idField              = SelectWhereColumn{"id"}
	kindField            = SelectWhereColumn{"kind"}
	parentKindField      = SelectWhereColumn{"parent_kind"}
	parentIdField        = SelectWhereColumn{"parent_id"}
	uniqueStringKeyField = WhereColumn{"unique_string_key"}
	stringKeyField       = WhereColumn{"string_key"}
	numericKeyField      = WhereColumn{"numeric_key"}
	timeKeyField         = WhereColumn{"time_key"}
	dataField            = SelectColumn{"data"}
	createTimeField      = SelectWhereColumn{"create_time"}
	updateTimeField      = SelectWhereColumn{"update_time"}
)

type PathField struct {
	p dialect.Path
}

func (p PathField) toWhereField() string {
	return p.p.JSONExtract("data")
}

func (p PathField) Asc() Order {
	return Order{p, false}
}

func (p PathField) Desc() Order {
	return Order{p, true}
}

func (p PathField) Key(key string) PathField {
	p.p.Key(key)
	return p
}

func (p PathField) Index(index int) PathField {
	p.p.Index(index)
	return p
}
