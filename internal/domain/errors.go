package domain

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("version conflict")
	ErrInvalid      = errors.New("invalid transition")
	ErrPrecondition = errors.New("precondition failed")
	ErrIdempotent   = errors.New("idempotent replay")
)
