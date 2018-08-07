package handlers

import (
	"github.com/andrskom/gorpcgen/protocol/context"
	"github.com/andrskom/gorpcgen/protocol/errors"
)

type GetRequest struct {
	Type string `json:"type"`
}

type GetResponse struct {
	Type string `json:"type"`
}

var HandlerErr = errors.New(errors.CodeBadRequest, "detail")
var HandlerOK = &GetResponse{Type: TypeOK}

const (
	TypeOK  = "ok"
	TypeErr = "err"
)

// goRpcGen:method
func (t *Test) Get(ctx *context.Context, req *GetRequest) (*GetResponse, *errors.Error) {
	switch req.Type {
	case TypeErr:
		return nil, HandlerErr
	case TypeOK:
		return HandlerOK, nil
	default:
		return nil, errors.New(errors.CodeInternal, "Unknown type of request")
	}
	return nil, nil
}
