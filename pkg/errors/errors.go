package errors

import "errors"

var (
	ErrBadReq		 = errors.New("bad request")
	ErrNotFound      = errors.New("not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrConflict      = errors.New("conflict")
	ErrInternalError = errors.New("internal error")
)