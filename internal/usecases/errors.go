package usecases

import "errors"

var (
	ErrNoRowsAffected = errors.New("no rows affected")
	ErrAlreadyExists  = errors.New("already exists")
	ErrNullFields     = errors.New("null field")
)
