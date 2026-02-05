package booking

import "errors"


var (
    ErrNotFound     = errors.New("booking not found")
    ErrForbidden    = errors.New("forbidden")
    ErrInvalidState = errors.New("invalid booking status")
	ErrDuplicateActiveRequest = errors.New("active request already exists")
    ErrConflict = errors.New("busy date")
)
