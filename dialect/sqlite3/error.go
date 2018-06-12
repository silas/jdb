package sqlite3

import (
	"github.com/mattn/go-sqlite3"
	"github.com/silas/jdb/internal/errors"
)

type sqlite3Error struct {
	err *sqlite3.Error
}

func newError(err *sqlite3.Error) errors.Error {
	return sqlite3Error{err: err}
}

func (e sqlite3Error) Error() string {
	return e.err.Error()
}

func (e sqlite3Error) Source() error {
	return e.err
}

func (e sqlite3Error) Type() errors.ErrorType {
	switch e.err.Code {
	case sqlite3.ErrPerm:
		return errors.PermissionError
	case sqlite3.ErrBusy, sqlite3.ErrLocked:
		return errors.BusyError
	case sqlite3.ErrConstraint:
		return errors.IntegrityError
	case sqlite3.ErrMismatch, sqlite3.ErrTooBig:
		return errors.DataError
	case sqlite3.ErrAuth:
		return errors.AuthorizationError
	default:
		return errors.UnknownError
	}
}
