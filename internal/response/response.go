package response

import (
	"fmt"
	"io"

	"github.com/JoStMc/httpfromtcp/internal/headers"
)

type StatusCode int
const (
	StatusOK 		  			StatusCode = 200
	StatusBadRequest  			StatusCode = 400
	StatusIntervalServerError	StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var out string
	switch statusCode {
	case StatusOK:
		out = fmt.Sprintf("HTTP/1.1 %d OK\r\n", statusCode)
	case StatusBadRequest:
		out = fmt.Sprintf("HTTP/1.1 %d Bad Request\r\n", statusCode)
	case StatusIntervalServerError:
		out = fmt.Sprintf("HTTP/1.1 %d Internal Server Error\r\n", statusCode)
	default:
		out = fmt.Sprintf("HTTP/1.1 %d\r\n", statusCode)
	}
	_, err := w.Write([]byte(out))
	return err
} 

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
} 

func WriteHeaders(w io.Writer, h headers.Headers) error {
	out := []byte{}
	h.ForEach(func(n, v string) {
		out = fmt.Appendf(out, "%s: %s\r\n", n, v)
	})
	out = fmt.Append(out, "\r\n")
	_, err := w.Write(out)
	return err
} 

