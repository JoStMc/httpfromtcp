package request

import (
	"errors"
	"io"
	"log"
	"strings"
)

type Request struct {
    RequestLine RequestLine
} 

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method		  string
} 

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	requestLine, err := parseRequestLine(string(data))
	if err != nil {
		return &Request{}, err
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
} 

func parseRequestLine(request string) (*RequestLine, error) {
	requestLine := strings.Split(request, "\r\n")[0]
	parts := strings.Fields(requestLine)
	if len(parts) != 3 {
	    return &RequestLine{}, errors.New("invalid request line")
	} 

	if parts[2][:5] != "HTTP/" {
	    return &RequestLine{}, errors.New("invalid version")
	} 
	version := strings.TrimPrefix(parts[2], "HTTP/")
	target := parts[1]
	method := parts[0]

	for _, ch := range method {
	    if ch < 'A' || ch > 'Z' {
	        return &RequestLine{}, errors.New("invalid method")
	    } 
	} 

	if version != "1.1" {
	    return &RequestLine{}, errors.New("HTTP version not supported")
	} 

	return &RequestLine{
		HttpVersion: version,
		RequestTarget: target,
		Method: method,
	}, nil
}
