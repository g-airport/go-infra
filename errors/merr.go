package errors

import (
	"encoding/json"

	errM "github.com/micro/go-micro/errors"
)

// micro server

func ParseRPCError(err error) *Error {
	if err == nil {
		return nil
	}

	errTxt := err.Error()

	// split "micro error" layer
	mErr := errM.Parse(errTxt)
	errTxt = mErr.Detail

	// try to parse to our service error
	return TryParse(errTxt)
}

func TryParse(err string) *Error {
	e := new(Error)
	if err := json.Unmarshal([]byte(err), e); err != nil {
		return nil
	}

	return e
}
