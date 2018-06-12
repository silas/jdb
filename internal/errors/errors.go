package errors

type ErrorType int

const (
	DataError ErrorType = iota
	IntegrityError
	TransactionError
	AuthorizationError
	BusyError
	PermissionError
	ConnectionError
	UnknownError
)

type Error interface {
	Error() string
	Source() error
	Type() ErrorType
}
