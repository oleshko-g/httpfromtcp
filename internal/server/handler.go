package server

import (
	"io"

	"github.com/oleshko-g/httpfromtcp/internal/request"
	"github.com/oleshko-g/httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (err HandlerError) writeError(w io.Writer) {
	response.WriteStatusLine(w, err.StatusCode)
	headers := response.GetDefaultHeaders(len(err.Message))
	response.WriteHeaders(w, headers)
	response.WriteBody(w, err.Message)
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerV2 func(w *response.Writer, req *request.Request)
