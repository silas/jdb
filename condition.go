package jdb

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	trueCondition  = "(1 = 1)"
	falseCondition = "(1 != 1)"
)

type expr struct {
	value string
}

func (c expr) toConditionSQL(query *bytes.Buffer, params *[]interface{}) error {
	query.WriteString(c.value)
	return nil
}

type Condition interface {
	toConditionSQL(sql *bytes.Buffer, params *[]interface{}) error
}

type conj []Condition

func (c conj) join(query *bytes.Buffer, params *[]interface{}, sep string) error {
	if len(c) > 0 {
		query.WriteString("(")
		for i, queryBuilder := range c {
			if i > 0 {
				query.WriteString(sep)
			}
			if queryBuilder == nil {
				return fmt.Errorf("%s: nil condition: %v", strings.TrimSpace(sep), c)
			}
			err := queryBuilder.toConditionSQL(query, params)
			if err != nil {
				return err
			}
		}
		query.WriteString(")")
	} else {
		query.WriteString(trueCondition)
	}
	return nil
}

type and conj

func (c and) toConditionSQL(query *bytes.Buffer, params *[]interface{}) error {
	return conj(c).join(query, params, " AND ")
}

func And(args ...Condition) Condition {
	return and(args)
}

type or conj

func (c or) toConditionSQL(query *bytes.Buffer, params *[]interface{}) error {
	return conj(c).join(query, params, " OR ")
}

func Or(args ...Condition) Condition {
	return or(args)
}

func True() Condition {
	return expr{trueCondition}
}

func False() Condition {
	return expr{falseCondition}
}

type eq struct {
	field WhereField
	value interface{}
}

func Eq(f WhereField, v interface{}) Condition {
	return eq{f, v}
}

func (c eq) toConditionSQL(query *bytes.Buffer, params *[]interface{}) error {
	if c.value != nil {
		query.WriteString(fmt.Sprintf("(%s = ?)", c.field.toWhereField()))
		*params = append(*params, c.value)
	} else {
		query.WriteString(fmt.Sprintf("(%s IS NULL)", c.field.toWhereField()))
	}
	return nil
}

type notEq struct {
	field WhereField
	value interface{}
}

func NotEq(f WhereField, v interface{}) Condition {
	return notEq{f, v}
}

func (c notEq) toConditionSQL(query *bytes.Buffer, params *[]interface{}) error {
	if c.value != nil {
		query.WriteString(fmt.Sprintf("(%s != ?)", c.field.toWhereField()))
		*params = append(*params, c.value)
	} else {
		query.WriteString(fmt.Sprintf("(%s IS NOT NULL)", c.field.toWhereField()))
	}
	return nil
}

type in struct {
	field WhereField
	value []interface{}
}

func In(f WhereField, v ...interface{}) Condition {
	return in{f, v}
}

func (c in) toConditionSQL(query *bytes.Buffer, params *[]interface{}) error {
	hasNil := false
	var value []interface{}
	for _, v := range c.value {
		if v != nil {
			value = append(value, v)
		} else {
			hasNil = true
		}
	}

	hasValues := len(value) > 0

	if !hasNil && !hasValues {
		query.WriteString(falseCondition)
		return nil
	}

	if hasNil && hasValues {
		query.WriteString("(")
	}
	if hasNil {
		if err := Eq(c.field, nil).toConditionSQL(query, params); err != nil {
			return err
		}
		if !hasValues {
			return nil
		}
	} else if !hasValues {
		query.WriteString(falseCondition)
		return nil
	}
	if hasNil && hasValues {
		query.WriteString(" OR ")
	}
	if hasValues {
		query.WriteString(fmt.Sprintf("(%s IN (%s))", c.field.toWhereField(), placeholders(len(value))))
		*params = append(*params, value...)
	}
	if hasNil && hasValues {
		query.WriteString(")")
	}

	return nil
}

type notIn struct {
	field WhereField
	value []interface{}
}

func NotIn(f WhereField, v ...interface{}) Condition {
	return notIn{f, v}
}

func (c notIn) toConditionSQL(query *bytes.Buffer, params *[]interface{}) error {
	hasNil := false
	var value []interface{}
	for _, v := range c.value {
		if v != nil {
			value = append(value, v)
		} else {
			hasNil = true
		}
	}

	hasValues := len(value) > 0

	if !hasNil && !hasValues {
		query.WriteString(trueCondition)
		return nil
	}

	if !hasNil && hasValues {
		query.WriteString("(")
	}
	if !hasNil {
		if err := Eq(c.field, nil).toConditionSQL(query, params); err != nil {
			return err
		}
	} else if !hasValues {
		if err := NotEq(c.field, nil).toConditionSQL(query, params); err != nil {
			return err
		}
	}
	if !hasNil && hasValues {
		query.WriteString(" OR ")
	}
	if hasValues {
		query.WriteString(fmt.Sprintf("(%s NOT IN (%s))", c.field.toWhereField(), placeholders(len(value))))
		*params = append(*params, value...)
	}
	if !hasNil && hasValues {
		query.WriteString(")")
	}

	return nil
}

type like struct {
	field WhereField
	value interface{}
}

func Like(f WhereField, v interface{}) Condition {
	return like{f, v}
}

func (c like) toConditionSQL(query *bytes.Buffer, a *[]interface{}) error {
	if c.value != nil {
		query.WriteString(fmt.Sprintf("(%s LIKE ?)", c.field.toWhereField()))
		*a = append(*a, c.value)
	} else {
		query.WriteString(fmt.Sprintf("(%s LIKE NULL)", c.field.toWhereField()))
	}
	return nil
}

type notLike struct {
	field WhereField
	value interface{}
}

func NotLike(f WhereField, v interface{}) Condition {
	return notLike{f, v}
}

func (c notLike) toConditionSQL(query *bytes.Buffer, params *[]interface{}) error {
	if c.value != nil {
		query.WriteString(fmt.Sprintf("(%s NOT LIKE ?)", c.field.toWhereField()))
		*params = append(*params, c.value)
	} else {
		query.WriteString(fmt.Sprintf("(%s NOT LIKE NULL)", c.field.toWhereField()))
	}
	return nil
}

type gt struct {
	field WhereField
	value interface{}
}

func Gt(f WhereField, v interface{}) Condition {
	return gt{f, v}
}

func (c gt) toConditionSQL(query *bytes.Buffer, params *[]interface{}) error {
	if c.value != nil {
		query.WriteString(fmt.Sprintf("(%s > ?)", c.field.toWhereField()))
		*params = append(*params, c.value)
	} else {
		query.WriteString(fmt.Sprintf("(%s > NULL)", c.field.toWhereField()))
	}
	return nil
}

type lt struct {
	field WhereField
	value interface{}
}

func Lt(f WhereField, v interface{}) Condition {
	return lt{f, v}
}

func (c lt) toConditionSQL(query *bytes.Buffer, params *[]interface{}) error {
	if c.value != nil {
		query.WriteString(fmt.Sprintf("(%s < ?)", c.field.toWhereField()))
		*params = append(*params, c.value)
	} else {
		query.WriteString(fmt.Sprintf("(%s < NULL)", c.field.toWhereField()))
	}
	return nil
}

type gte struct {
	field WhereField
	value interface{}
}

func Gte(f WhereField, v interface{}) Condition {
	return gte{f, v}
}

func (c gte) toConditionSQL(query *bytes.Buffer, params *[]interface{}) error {
	if c.value != nil {
		query.WriteString(fmt.Sprintf("(%s >= ?)", c.field.toWhereField()))
		*params = append(*params, c.value)
	} else {
		query.WriteString(fmt.Sprintf("(%s >= NULL)", c.field.toWhereField()))
	}
	return nil
}

type lte struct {
	field WhereField
	value interface{}
}

func Lte(f WhereField, v interface{}) Condition {
	return lte{f, v}
}

func (c lte) toConditionSQL(query *bytes.Buffer, params *[]interface{}) error {
	if c.value != nil {
		query.WriteString(fmt.Sprintf("(%s <= ?)", c.field.toWhereField()))
		*params = append(*params, c.value)
	} else {
		query.WriteString(fmt.Sprintf("(%s <= NULL)", c.field.toWhereField()))
	}
	return nil
}
