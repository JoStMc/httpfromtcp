package request

import (
	"errors"
	"io"
	"strings"
)

var SEPARATOR = "\r\n"
var bufferSize = 8

type state int
const (
	initialized state = iota
	done
)

type Request struct {
    RequestLine RequestLine
	State		state
} 

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method		  string
} 


func newRequest() *Request {
    return &Request{
		State: state(initialized),
	}
} 

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest() 

	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0

	for request.State != state(done) {
		if buf[len(buf)-1] != 0 {
			newBuf := make([]byte, cap(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		} 

		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.State = state(done);
				break
			} 
			return nil, err
		}
		readToIndex += n

		bytesParsed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[bytesParsed:readToIndex])
		readToIndex -= bytesParsed
	} 

	return request, nil
} 

func parseRequestLine(b string) (*RequestLine, int, error) {
	idx := strings.Index(b, SEPARATOR)
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
	}, idx, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.State == state(initialized) {
		requestLine, bytesParsed, err := parseRequestLine(string(data))
		if err != nil {
			return bytesParsed, nil
		}
		if bytesParsed == 0 {
		    return 0, nil
		} 
		r.RequestLine = *requestLine
		r.State++
		return bytesParsed, nil
	} else if r.State == state(done) {
		return 0, errors.New("error: trying to read data in a done state")
	} else {
		return 0, errors.New("error: unknown state")
	} 
} 
