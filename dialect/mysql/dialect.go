package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/silas/jdb"
	"github.com/silas/jdb/dialect"
)

const driverName = `mysql`
const timestamp = `NOW(4)`

type mysqlDialect struct{}

func init() {
	jdb.RegisterDialect(driverName, &mysqlDialect{})
}

func (d *mysqlDialect) OrderExpression(o dialect.OrderField) string {
	if o.OrderDesc() {
		return fmt.Sprintf("%s DESC", o.OrderField())
	} else {
		return fmt.Sprintf("%s ASC", o.OrderField())
	}
}

func (d *mysqlDialect) ReplacePlaceHolders(text string) string {
	return text
}

func (d *mysqlDialect) TimestampExpression() string {
	return timestamp
}

func (d *mysqlDialect) Migrate(ctx context.Context, db *sql.DB, table string) error {
	return revisions.Run(ctx, db, &migrationHelper{table})
}

func (d *mysqlDialect) ValidateDataSourceName(dsn string, opts dialect.ValidateDataSourceNameOpts) error {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return err
	}

	if !cfg.ParseTime {
		return errors.New("parseTime is required")
	}

	return nil
}

func (d *mysqlDialect) Path() dialect.Path {
	return &mysqlPath{}
}

func (d *mysqlDialect) Now(ctx context.Context, tx *sql.Tx) (now time.Time, err error) {
	query := fmt.Sprintf("SELECT %s AS now", timestamp)
	row := tx.QueryRowContext(ctx, query)
	err = row.Scan(&now)
	return
}

func (d *mysqlDialect) ErrorMap(err error) error {
	if e, ok := err.(*mysql.MySQLError); ok {
		return newError(e)
	}
	return err
}
