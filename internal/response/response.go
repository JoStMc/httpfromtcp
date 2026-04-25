package response

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/JoStMc/httpfromtcp/internal/headers"
)

var separator = []byte("\r\n")

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
	if w.writerState != writerState(stateHeaders) {
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
	if w.writerState != stateBody {
	    return 0, responseOutOfOrder
	} 
	lenSep := len(separator)

	idx := bytes.Index(p, separator)
	if idx == -1 {
	    return 0, nil
	} 
	idx2 := bytes.Index(p[lenSep:], separator)
	if idx2 == -1 {
	    return 0, nil
	} 
	
	strChunkLen := string(p[:idx])
	// may be a bad way to do this
	chunkLen, err := strconv.ParseInt(strChunkLen, 16, len(strChunkLen)*4)
	if err != nil {
		return 0, err
	}

	if idx2 != int(chunkLen) {
		return 0, errors.New("chunk length mismatch")
	} 

	if chunkLen == 0 {
		// bytes parsed
		return len(strChunkLen) + len(separator)*2, nil
	} 

	chunk := p[idx+lenSep:idx+lenSep+idx2]

	_, err = w.writer.Write(chunk)
	if err != nil {
		return 0, err
	}

	bytesParsed := idx+2*lenSep+idx
	n, err := w.WriteChunkedBody(p[bytesParsed:])
	bytesParsed+=n
	return bytesParsed, err
}
