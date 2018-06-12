package jdb

import (
	"errors"

	jdberrors "github.com/silas/jdb/internal/errors"
)

type Error = jdberrors.Error
type ErrorType = jdberrors.ErrorType

var (
	ErrReadOnlyMode = errors.New("jdb: read-only mode")
	ErrIDNotFound   = errors.New("jdb: id not found")
	ErrNotFound     = errors.New("jdb: not found")
)
