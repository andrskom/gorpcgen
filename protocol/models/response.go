package models

import (
	"encoding/json"
	"github.com/andrskom/gorpcgen/protocol/errors"
)

// Response is response object.
// Contains RespMeta field for meta information
// U can change default response object in templates.
type Response struct {
	Meta   RespMeta        `json:"meta"`
	Type   ResponseType    `json:"type"`
	Result json.RawMessage `json:"result"`
}

// RepType is type of response(err or ok)
type ResponseType string

const (
	// ResponseTypeOK is type for positive result
	ResponseTypeOK ResponseType = "ok"
	// ResponseTypeErr is type for err result
	ResponseTypeErr ResponseType = "err"
)

// RespMeta is meta information for response.
// That information can be add to logging for example.
type RespMeta struct {
	RequestID string `json:"request_id"`
}

func NewJSONMarshallErrResponseBytes(requestID string) []byte {
	return []byte(
		`{
	"meta": {
		"request_id": "` + requestID + `"
	},
	"type", "` + string(ResponseTypeErr) + `",
	"result": {
		"code": "` + string(errors.CodeInternal) + `",
		"msg": "Can't marshal response'"
		}
	}`,
	)
}

func NewResponseErr(requestID string, error *errors.Error) *Response {
	bytes, err := json.Marshal(error)
	if err != nil {
		return &Response{
			Meta: RespMeta{RequestID: requestID},
			Type: ResponseTypeErr,
			Result: json.RawMessage(
				`{
	"code": "` + errors.CodeInternal + `",
	"attributes": null,
	"detail": "Can't marshal error to response"
	}`,
			),
		}
	}
	return &Response{Meta: RespMeta{RequestID: requestID}, Type: ResponseTypeErr, Result: json.RawMessage(bytes)}
}

func NewResponseOK(requestID string, data []byte) *Response {
	return &Response{Meta: RespMeta{RequestID: requestID}, Type: ResponseTypeOK, Result: json.RawMessage(data)}
}
