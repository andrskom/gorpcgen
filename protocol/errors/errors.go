package errors

// Error is error type of RPC error, u can use custom object, simple edit template.
// Use code for build your logic on error(if u like it).
type Error struct {
	Code       Code              `json:"code"`
	Detail     string            `json:"detail"`
	Attributes map[string]string `json:"attributes"`
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
)
