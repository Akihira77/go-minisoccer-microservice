package error

import "errors"

const (
	Success = "success"
	Error   = "error"
)

var (
	ErrInternalServer = errors.New("internal server error")
	ErrSQL            = errors.New("database server failed to execute query")
	ErrTooManyRequest = errors.New("too many requests")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrInvalidToken   = errors.New("invalid token")
	ErrForbidden      = errors.New("forbidden")
)

var GeneralErrors = []error{
	ErrInternalServer,
	ErrSQL,
	ErrTooManyRequest,
	ErrUnauthorized,
	ErrInvalidToken,
	ErrForbidden,
}
