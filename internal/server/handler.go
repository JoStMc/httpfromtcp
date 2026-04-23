package server

import (
	"io"

	"github.com/JoStMc/httpfromtcp/internal/request"
	"github.com/JoStMc/httpfromtcp/internal/response"
)

type HandlerError struct {
	statusCode   response.StatusCode
	errorMessage []byte
} 

type Handler func(w io.Writer, req *request.Request) *HandlerError


func NewHandlerError(code response.StatusCode, message []byte) *HandlerError {
    return &HandlerError{
		statusCode: code,
		errorMessage: message,
	}
} 
