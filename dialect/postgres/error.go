package postgres

import (
	"github.com/lib/pq"
	"github.com/silas/jdb/internal/errors"
)

type postgresError struct {
	err *pq.Error
}

func newError(err *pq.Error) errors.Error {
	return postgresError{err: err}
}

func (e postgresError) Error() string {
	return e.err.Error()
}

func (e postgresError) Source() error {
	return e.err
}

func (e postgresError) Type() errors.ErrorType {
	switch e.err.Code.Class() {
	case "08":
		return errors.ConnectionError
	case "22":
		return errors.DataError
	case "23":
		return errors.IntegrityError
	case "25":
		return errors.TransactionError
	case "28":
		return errors.AuthorizationError
	default:
		return errors.UnknownError
	}
}
