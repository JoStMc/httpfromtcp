package server

import (
	"github.com/JoStMc/httpfromtcp/internal/request"
	"github.com/JoStMc/httpfromtcp/internal/response"
)

type HandlerError struct {
	statusCode   response.StatusCode
	errorMessage []byte
} 

type Handler func(w *response.Writer, req *request.Request)

func NewHandlerError(code response.StatusCode, message []byte) *HandlerError {
    return &HandlerError{
		statusCode: code,
		errorMessage: message,
	}
} 
