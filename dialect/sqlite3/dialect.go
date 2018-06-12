package sqlite3

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/silas/jdb"
	"github.com/silas/jdb/dialect"
)

const driverName = `sqlite3`
const timestamp = `STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')`
const layout = `2006-01-02 15:04:05.000`

type sqlite3Dialect struct{}

func init() {
	jdb.RegisterDialect(driverName, &sqlite3Dialect{})
}

func (d *sqlite3Dialect) OrderExpression(o dialect.OrderField) string {
	if o.OrderDesc() {
		return fmt.Sprintf("%s DESC", o.OrderField())
	} else {
		return fmt.Sprintf("%s ASC", o.OrderField())
	}
}

func (d *sqlite3Dialect) ReplacePlaceHolders(text string) string {
	return text
}

func (d *sqlite3Dialect) TimestampExpression() string {
	return timestamp
}

func (d *sqlite3Dialect) Migrate(ctx context.Context, db *sql.DB, table string) error {
	return revisions.Run(ctx, db, &migrationHelper{table})
}

func (d *sqlite3Dialect) ValidateDataSourceName(dsn string, opts dialect.ValidateDataSourceNameOpts) error {
	readOnly := false
	pos := strings.IndexRune(dsn, '?')
	if pos >= 1 {
		params, err := url.ParseQuery(dsn[pos+1:])
		if err != nil {
			return err
		}

		mode := params.Get("mode")
		if mode == "ro" {
			readOnly = true
		}
	}

	if opts.ReadOnly && !readOnly {
		return errors.New("expected mode=ro")
	} else if !opts.ReadOnly && readOnly {
		return errors.New("unexpected mode=ro")
	}

	return nil
}

func (d *sqlite3Dialect) Path() dialect.Path {
	return &sqlite3Path{}
}

func (d *sqlite3Dialect) Now(ctx context.Context, tx *sql.Tx) (now time.Time, err error) {
	query := fmt.Sprintf("SELECT %s AS now", timestamp)
	row := tx.QueryRowContext(ctx, query)
	err = row.Scan((*nowTime)(&now))
	return
}

type nowTime time.Time

func (t *nowTime) Scan(v interface{}) error {
	vt, err := time.Parse(layout, string(v.([]byte)))
	if err != nil {
		return err
	}
	*t = nowTime(vt)
	return nil
}

func (d *sqlite3Dialect) ErrorMap(err error) error {
	if e, ok := err.(*sqlite3.Error); ok {
		return newError(e)
	}
	return err
}
