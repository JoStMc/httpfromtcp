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


type writerState int
const (
	stateStatusLine writerState = iota
	stateHeaders
	stateBody
	stateTrailers
)
var responseOutOfOrder = fmt.Errorf("response not written in the correct order")

type Writer struct {
	writer io.Writer
	writerState writerState
} 

func NewWriter(writer io.Writer) *Writer {
	return &Writer{writer: writer, writerState: stateStatusLine}
} 

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != writerState(stateStatusLine) {
	    return responseOutOfOrder
	} 
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
	_, err := w.writer.Write([]byte(out))
	w.writerState++
	return err
} 

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.writerState != stateHeaders && w.writerState != stateTrailers {
	    return responseOutOfOrder
	} 
	out := []byte{}
	h.ForEach(func(n, v string) {
		out = fmt.Appendf(out, "%s: %s\r\n", n, v)
	})
	out = fmt.Append(out, "\r\n")
	_, err := w.writer.Write(out)
	w.writerState++
	return err
} 

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writerState(stateBody) {
	    return 0, responseOutOfOrder
	} 
	return w.writer.Write(p)
} 

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
} 


func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	length := len(p)
	w.WriteBody(fmt.Appendf(nil, "%x\r\n", length))
	n, err := w.WriteBody(p)
	w.WriteBody([]byte("\r\n"))
	return n, err}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	out := []byte("0\r\n")
	w.writerState++
	return w.writer.Write(out)
} 

func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.writerState != stateTrailers {
	    return responseOutOfOrder
	} 
	err := w.WriteHeaders(h)
	return err
} 
