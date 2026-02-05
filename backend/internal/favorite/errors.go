package favorite

import "errors"

var (
	ErrAlreadyExists = errors.New("favorite already exists")
	ErrNotFound      = errors.New("favorite not found")
	ErrForbidden     = errors.New("forbidden")
)
