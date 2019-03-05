package auth

import (
	"encoding/json"

	"github.com/go-infra/errors"

	errmicro "github.com/micro/go-micro/errors"
)

var (
	ErrNoMethod     = errors.NotFound(40001, "no method")
	ErrNoToken      = errors.Unauthorized(40002, "no token")
	ErrNoPermission = errors.Forbidden(40003, "no permission")
)

func ParseFromError(err error) *errors.Error {
	if err == nil {
		return nil
	}

	errTxt := err.Error()

	// split "micro error" layer
	mErr := errmicro.Parse(errTxt)
	errTxt = mErr.Detail

	// try to parse to our service error
	return TryParse(errTxt)
}

func TryParse(err string) *errors.Error {
	e := new(errors.Error)
	if err := json.Unmarshal([]byte(err), e); err != nil {
		return nil
	}

	return e
}
