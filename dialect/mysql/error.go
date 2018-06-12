package mysql

import (
	"github.com/go-sql-driver/mysql"
	"github.com/silas/jdb/internal/errors"
)

type postgresError struct {
	err *mysql.MySQLError
}

func newError(err *mysql.MySQLError) errors.Error {
	return postgresError{err: err}
}

func (e postgresError) Error() string {
	return e.err.Error()
}

func (e postgresError) Source() error {
	return e.err
}

func (e postgresError) Type() errors.ErrorType {
	switch e.err.Number {
	case 1205:
		return errors.BusyError
	case 1213:
		return errors.TransactionError
	default:
		return errors.UnknownError
	}
}
