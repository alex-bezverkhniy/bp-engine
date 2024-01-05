package engine

import "errors"

var (
	ErrProcessNotFound     error = errors.New("process not found")
	ErrCannotCreateProcess error = errors.New("cannot create process")
)
