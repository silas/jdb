package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"

	"time"

	"github.com/lib/pq"
	"github.com/silas/jdb"
	"github.com/silas/jdb/dialect"
)

const driverName = `postgres`
const timestamp = `CURRENT_TIMESTAMP`

var placeHolder = regexp.MustCompile(`\?`)

type postgresDialect struct{}

func init() {
	jdb.RegisterDialect(driverName, &postgresDialect{})
}

func (d *postgresDialect) OrderExpression(o dialect.OrderField) string {
	if o.OrderDesc() {
		return fmt.Sprintf("%s DESC NULLS LAST", o.OrderField())
	} else {
		return fmt.Sprintf("%s ASC NULLS FIRST", o.OrderField())
	}
}

func (d *postgresDialect) ReplacePlaceHolders(text string) string {
	i := 0
	dollarInc := func(b []byte) []byte {
		i++
		return []byte(fmt.Sprintf("$%d", i))
	}
	return string(placeHolder.ReplaceAllFunc([]byte(text), dollarInc))
}

func (d *postgresDialect) TimestampExpression() string {
	return timestamp
}

func (d *postgresDialect) Migrate(ctx context.Context, db *sql.DB, table string) error {
	return revisions.Run(ctx, db, &migrationHelper{table})
}

func (d *postgresDialect) ValidateDataSourceName(dsn string, opts dialect.ValidateDataSourceNameOpts) error {
	return nil
}

func (d *postgresDialect) Path() dialect.Path {
	return &postgresPath{}
}

func (d *postgresDialect) Now(ctx context.Context, tx *sql.Tx) (now time.Time, err error) {
	query := fmt.Sprintf("SELECT %s AS now", timestamp)
	row := tx.QueryRowContext(ctx, query)
	err = row.Scan(&now)
	return
}

func (d *postgresDialect) ErrorMap(err error) error {
	if e, ok := err.(*pq.Error); ok {
		return newError(e)
	}
	return err
}
