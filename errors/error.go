package errors

import (
	"encoding/json"
	"log"
)

//----------------------------------------------------*
// recommend : github.com/pkg/errors
//----------------------------------------------------*

type Error struct {
	Code     int         `json:"code"`
	Status   int         `json:"Status"`
	Detail   string      `json:"detail"`
	Internal string      `json:"internal,omitempty"`
	Content  interface{} `json:"content,omitempty"`
}

func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// map[err code]error
var Errors = map[int]*Error{}

func Internal(detail string, err error) error {

	internal := ""
	if err != nil {
		internal = err.Error()
	}

	return &Error{
		Status:   500, // 🤣 use internal service err
		Detail:   detail,
		Internal: internal,
	}
}

//----------------------------------------------------*
// base http code err
//----------------------------------------------------*

func addError(err *Error) *Error {
	e, ok := Errors[err.Code]
	if ok {
		log.Fatalf("duplate error code: %v, %v", e, err)
	}

	Errors[err.Code] = err
	return err
}

func BadRequest(code int, detail string) error {
	return addError(&Error{
		Code:   code,
		Status: 400,
		Detail: detail,
	})
}

func Conflict(code int, detail string) error {
	return addError(&Error{
		Code:   code,
		Status: 409,
		Detail: detail,
	})
}

func Unauthorized(code int, detail string) error {
	return addError(&Error{
		Code:   code,
		Status: 401,
		Detail: detail,
	})
}

func Forbidden(code int, detail string) error {
	return addError(&Error{
		Code:   code,
		Status: 403,
		Detail: detail,
	})
}

func NotFound(code int, detail string) error {
	return addError(&Error{
		Code:   code,
		Status: 404,
		Detail: detail,
	})
}
