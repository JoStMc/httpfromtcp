package request

import (
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/JoStMc/httpfromtcp/internal/headers"
)

var separator = "\r\n"
var bufferSize = 8

type state int
const (
	requestStateInitialized state = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

type Request struct {
    RequestLine RequestLine
	Headers 	headers.Headers
	State		state
	Body 		[]byte
} 

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method		  string
} 


func newRequest() *Request {
    return &Request{
		State: requestStateInitialized,
		Headers: headers.NewHeaders(),
	}
} 

func (r *Request) GetBody() string {
    return string(r.Body)
} 

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest() 

	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	reachedEnd := false

	for request.State != requestStateDone {
		if buf[len(buf)-1] != 0 {
			newBuf := make([]byte, cap(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		} 

		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				reachedEnd = true
			} else {
				return nil, err
			} 
		}
		readToIndex += n

		bytesParsed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		if reachedEnd {
		    request.State = requestStateDone
			break
		} 
		copy(buf, buf[bytesParsed:readToIndex])
		readToIndex -= bytesParsed
	} 

	// Not a fan of this, but apparently Content-Length > actual body isn't
	// deseriable and there's no easy way to check we have EOF in .parse
	if cl, _ := request.getContentLength(); len(request.Body) != cl {
		return nil, errors.New("body shorter than expected")
	}
	return request, nil
} 

func parseRequestLine(b string) (*RequestLine, int, error) {
	idx := strings.Index(b, separator)
	if idx == -1 {
	    return nil, 0, nil
	} 

	requestLine := b[:idx]
	parts := strings.Fields(requestLine)
	if len(parts) != 3 {
	    return nil, idx, errors.New("invalid request line")
	} 

	if parts[2][:5] != "HTTP/" {
	    return nil, idx, errors.New("invalid version")
	} 
	version := strings.TrimPrefix(parts[2], "HTTP/")
	target := parts[1]
	method := parts[0]

	for _, ch := range method {
	    if ch < 'A' || ch > 'Z' {
	        return nil, idx, errors.New("invalid method")
	    } 
	} 

	if version != "1.1" {
	    return nil, idx, errors.New("HTTP version not supported")
	} 

	return &RequestLine{
		HttpVersion: version,
		RequestTarget: target,
		Method: method,
	}, idx+len(separator), nil
}

// Could make more general to int headers
func (r *Request) getContentLength() (int, error) {
	contentLength := r.Headers.Get("content-length")
	if contentLength == "" {
		return 0, nil
	}
	length, err := strconv.Atoi(contentLength)
	if err != nil {
		return 0, errors.New("error: invalid content-length header")
	}
	return length, nil
} 

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case requestStateInitialized:
		requestLine, bytesParsed, err := parseRequestLine(string(data))
		if err != nil {
			return bytesParsed, err
		}
		if bytesParsed == 0 {
			return 0, nil
		} 
		r.RequestLine = *requestLine
		r.State++
		return bytesParsed, nil
	case requestStateParsingHeaders:
		bytesParsed, done, err := r.Headers.Parse(data)
		if err != nil {
			return bytesParsed, err
		}
		if done {
			r.State++
		} 
		return bytesParsed, nil
	case requestStateParsingBody:
		length, err := r.getContentLength()
		if err != nil {
			return 0, err
		}
		if length == 0 {
		    r.State++
			return 0, nil
		} 

		r.Body = append(r.Body, data...) 

		if len(r.Body) > length {
			return len(data), errors.New("error: body longer than expected")
		} else if len(r.Body) == length {
			r.State++
		} 
		return len(data), nil
	case requestStateDone:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("error: unknown state")
	}
} 
