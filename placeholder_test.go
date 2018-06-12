package jdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlaceholders(t *testing.T) {
	require.Equal(t, placeholders(0), "")
	require.Equal(t, placeholders(1), "?")
	require.Equal(t, placeholders(2), "?, ?")
	require.Equal(t, placeholders(3), "?, ?, ?")
}
