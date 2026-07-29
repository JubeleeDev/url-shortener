package shortener

import "errors"

var (
	ErrInvalidURL        = errors.New("invalid url")
	ErrInvalidCodeLength = errors.New("invalid code length")
	ErrNotFound          = errors.New("not found")
	ErrUniqueConflict    = errors.New("unique code conflict")
)
