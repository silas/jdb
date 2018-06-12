package dialect

import (
	"context"
	"database/sql"
	"time"
)

type Dialect interface {
	ValidateDataSourceName(v string, opts ValidateDataSourceNameOpts) error

	OrderExpression(order OrderField) string
	ReplacePlaceHolders(sql string) string
	TimestampExpression() string
	Path() Path
	Now(ctx context.Context, tx *sql.Tx) (time.Time, error)
	ErrorMap(err error) error

	Migrate(ctx context.Context, db *sql.DB, table string) error
}

type ValidateDataSourceNameOpts struct {
	ReadOnly bool
}
