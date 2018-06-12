package jdb

import (
	"bytes"
	"testing"

	"fmt"

	"github.com/stretchr/testify/require"
)

func TestConditions(t *testing.T) {
	tests := []struct {
		Condition Condition
		Query     string
		Params    []interface{}
	}{
		// True
		{
			True(),
			"(1 = 1)",
			params(),
		},
		// False
		{
			False(),
			"(1 != 1)",
			params(),
		},
		// Eq
		{
			Eq(idField, "1"),
			"(id = ?)",
			params("1"),
		},
		{
			Eq(idField, nil),
			"(id IS NULL)",
			params(),
		},
		// NotEq
		{
			NotEq(idField, "1"),
			"(id != ?)",
			params("1"),
		},
		{
			NotEq(idField, nil),
			"(id IS NOT NULL)",
			params(),
		},
		// In
		{
			In(idField),
			falseCondition,
			params(),
		},
		{
			In(idField, nil),
			"(id IS NULL)",
			params(),
		},
		{
			In(idField, "1"),
			"(id IN (?))",
			params("1"),
		},
		{
			In(idField, "1", "2"),
			"(id IN (?, ?))",
			params("1", "2"),
		},
		{
			In(idField, nil, "1", "2"),
			"((id IS NULL) OR (id IN (?, ?)))",
			params("1", "2"),
		},
		// NotIn
		{
			NotIn(idField),
			trueCondition,
			params(),
		},
		{
			NotIn(idField, nil),
			"(id IS NOT NULL)",
			params(),
		},
		{
			NotIn(idField, "1"),
			"((id IS NULL) OR (id NOT IN (?)))",
			params("1"),
		},
		{
			NotIn(idField, "1", "2"),
			"((id IS NULL) OR (id NOT IN (?, ?)))",
			params("1", "2"),
		},
		{
			NotIn(idField, nil, "1", "2"),
			"(id NOT IN (?, ?))",
			params("1", "2"),
		},
		// Like
		{
			Like(stringKeyField, "john@%"),
			"(string_key LIKE ?)",
			params("john@%"),
		},
		{
			Like(stringKeyField, nil),
			"(string_key LIKE NULL)",
			params(),
		},
		// NotLike
		{
			NotLike(stringKeyField, "john@%"),
			"(string_key NOT LIKE ?)",
			params("john@%"),
		},
		{
			NotLike(stringKeyField, nil),
			"(string_key NOT LIKE NULL)",
			params(),
		},
		// Gt
		{
			Gt(numericKeyField, 10),
			"(numeric_key > ?)",
			params(10),
		},
		{
			Gt(numericKeyField, nil),
			"(numeric_key > NULL)",
			params(),
		},
		// Lt
		{
			Lt(numericKeyField, 10),
			"(numeric_key < ?)",
			params(10),
		},
		{
			Lt(numericKeyField, nil),
			"(numeric_key < NULL)",
			params(),
		},
		// Gte
		{
			Gte(numericKeyField, 10),
			"(numeric_key >= ?)",
			params(10),
		},
		{
			Gte(numericKeyField, nil),
			"(numeric_key >= NULL)",
			params(),
		},
		// Lte
		{
			Lte(numericKeyField, 10),
			"(numeric_key <= ?)",
			params(10),
		},
		{
			Lte(numericKeyField, nil),
			"(numeric_key <= NULL)",
			params(),
		},
		// And
		{
			And(),
			"(1 = 1)",
			params(),
		},
		{
			And(Eq(kindField, "test")),
			"((kind = ?))",
			params("test"),
		},
		{
			And(Eq(kindField, "test"), Eq(idField, "1")),
			"((kind = ?) AND (id = ?))",
			params("test", "1"),
		},
		// Or
		{
			Or(),
			"(1 = 1)",
			params(),
		},
		{
			Or(Eq(kindField, "test")),
			"((kind = ?))",
			params("test"),
		},
		{
			Or(Eq(kindField, "test"), Eq(idField, "1")),
			"((kind = ?) OR (id = ?))",
			params("test", "1"),
		},
	}

	for i, test := range tests {
		msg := fmt.Sprintf("Test: %d", i)

		var params []interface{}
		query := &bytes.Buffer{}

		err := test.Condition.toConditionSQL(query, &params)
		require.NoError(t, err, msg)
		require.Equal(t, test.Query, query.String(), msg)
		require.Equal(t, test.Params, params, msg)
	}
}
