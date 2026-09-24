package repository

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrDuplicateCode = errors.New("short code already in use")
)
