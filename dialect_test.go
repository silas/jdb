package jdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDialect(t *testing.T) {
	_, err := Dialect("unknown")
	require.EqualError(t, err, `unknown dialect "unknown" (forgotten import?)`)

	d, err := Dialect("sqlmock")
	require.NoError(t, err)
	require.NotNil(t, d)
}
