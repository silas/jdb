package sqlmock

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"time"

	"github.com/silas/jdb/dialect"
	_ "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const driverName = `sqlmock`
const timestamp = `CURRENT_TIMESTAMP`

type mockDialect struct{}

func RegisterDialectArgs() (string, dialect.Dialect) {
	return driverName, &mockDialect{}
}

func (d *mockDialect) OrderExpression(o dialect.OrderField) string {
	if o.OrderDesc() {
		return fmt.Sprintf("%s DESC", o.OrderField())
	} else {
		return fmt.Sprintf("%s ASC", o.OrderField())
	}
}

func (d *mockDialect) ReplacePlaceHolders(text string) string {
	return text
}

func (d *mockDialect) TimestampExpression() string {
	return timestamp
}

func (d *mockDialect) Migrate(ctx context.Context, db *sql.DB, table string) error {
	return revisions.Run(ctx, db, &migrationHelper{table})
}

func (d *mockDialect) ValidateDataSourceName(dsn string, opts dialect.ValidateDataSourceNameOpts) error {
	if strings.Contains(dsn, "require-read-only-true") && !opts.ReadOnly {
		return errors.New("ReadOnly=false")
	}
	if strings.Contains(dsn, "require-read-only-false") && opts.ReadOnly {
		return errors.New("ReadOnly=true")
	}

	return nil
}

func (d *mockDialect) Path() dialect.Path {
	return &mockPath{}
}

func (d *mockDialect) Now(ctx context.Context, tx *sql.Tx) (now time.Time, err error) {
	query := fmt.Sprintf("SELECT %s AS now", timestamp)
	row := tx.QueryRowContext(ctx, query)
	err = row.Scan(&now)
	return
}

func (d *mockDialect) ErrorMap(err error) error {
	return err
}
