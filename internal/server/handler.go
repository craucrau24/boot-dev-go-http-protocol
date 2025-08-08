package server

import (
	"io"

	"github.com/craucrau24/boot-dev-go-http-protocol/internal/request"
	"github.com/craucrau24/boot-dev-go-http-protocol/internal/response"
)

type HandlerError struct {
	Status response.StatusCode
	Message string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError