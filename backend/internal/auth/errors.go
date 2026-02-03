package auth

import "errors"

var (
    ErrInvalidRefresh     = errors.New("invalid refresh token")
    ErrRefreshReuse       = errors.New("refresh token reuse detected")
    ErrRefreshGraceReuse  = errors.New("refresh token reuse within grace period")
    ErrRefreshAlreadyRotated = errors.New("refresh already rotated recently")
)
