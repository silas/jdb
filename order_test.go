package jdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrder_OrderField(t *testing.T) {
	o := Order{idField, true}
	require.Equal(t, o.OrderField(), idField.toWhereField())
}

func TestOrder_OrderDesc(t *testing.T) {
	o := Order{idField, true}
	require.True(t, o.OrderDesc())
}
