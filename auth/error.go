package auth

import (
	iErr "github.com/go-infra/errors"
)

var (
	ErrNoMethod     = iErr.NotFound(40001, "no method")
	ErrNoToken      = iErr.Unauthorized(40002, "no token")
	ErrNoPermission = iErr.Forbidden(40003, "no permission")
)
