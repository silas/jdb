package jdb

import (
	"fmt"

	"github.com/silas/jdb/dialect"
)

var dialects = make(map[string]dialect.Dialect)

func Dialect(driverName string) (dialect.Dialect, error) {
	if d, ok := dialects[driverName]; ok {
		return d, nil
	} else {
		return nil, fmt.Errorf(`unknown dialect "%s" (forgotten import?)`, driverName)
	}
}

func RegisterDialect(driverName string, fn dialect.Dialect) {
	if _, exists := dialects[driverName]; exists {
		panic("duplicate dialect registration " + driverName)
	}
	dialects[driverName] = fn
}
