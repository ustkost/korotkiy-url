package service

import "errors"

const maxListLimit = 100

var (
	ErrInvalidLimit  = errors.New("limit must be between 1 and 100")
	ErrInvalidOffset = errors.New("offset must be non-negative")
)
