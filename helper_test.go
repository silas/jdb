package jdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func params(v ...interface{}) []interface{} {
	if len(v) != 0 {
		return v
	}
	return nil
}

func setupQuery(t *testing.T) *Query {
	d, err := Dialect("sqlmock")
	require.NoError(t, err)

	return newQuery(d, "table", "kind")
}
