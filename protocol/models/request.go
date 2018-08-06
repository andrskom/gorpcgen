package models

import "encoding/json"

// Request is base request object
// Contains RequestMeta field for meta information
// U can change default request object in templates.
type Request struct {
	Meta RequestMeta     `json:"meta"`
	Data json.RawMessage `json:"data"`
}

// RequestMeta is meta information for request.
// That information can be add to logging for example.
type RequestMeta struct {
	RequestID string `json:"request_id"`
	Origin    string `json:"origin"`
}
