package errors

var (
	ErrInvalidPing      = BadRequest(50001, "bad ping")
	ErrDecodeArgsFailed = BadRequest(50002, "decode args failed")
)
