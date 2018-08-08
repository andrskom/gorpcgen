package errors

// Error is error type of RPC error, u can use custom object, simple edit template.
// Use code for build your logic on error(if u like it).
type Error struct {
	Code   Code     `json:"code"`
	Detail string   `json:"detail"`
	Errors []string `json:"errors"`
}

// Code is type for encapsulating Codes
type Code string

const (
	// CodeInternal for internal server error
	CodeInternal Code = "internal"
	// CodeBadRequest for bad request from client
	CodeBadRequest Code = "bad_request"
	// CodeNotFound if u can't find something
	CodeNotFound Code = "not_found"
	// CodeRestricted if action is restricted for user
	CodeRestricted Code = "restricted"
	// CodeProtocol for protocol specified errors
	CodeProtocol Code = "protocol"
	// CodeClient for client specified errors
	CodeClient Code = "client"
)

func New(code Code, detail string, errs ...error) *Error {
	e := &Error{
		Code:   code,
		Detail: detail,
		Errors: make([]string, len(errs)),
	}
	for i := 0; i < len(errs); i++ {
		e.Errors[i] = errs[i].Error()
	}
	return e
}
