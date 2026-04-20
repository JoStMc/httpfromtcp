package headers

import (
	"bytes"
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
} 

var separator = []byte("\r\n")

func parseHeaderLine(line []byte) (string, string, error) {
	parts := bytes.Fields(line)
	if len(parts) != 2 {
		return "", "", errors.New("invalid header line")
	} 
	field_name := string(parts[0])
	field_value := string(parts[1])
	if field_name[len(field_name)-1] != ':' {
		return "", "", errors.New("invalid header line")
	} 
	return strings.TrimRight(field_name, ":"), field_value, nil
} 

func (h Headers) Parse(data []byte) (int, bool, error) {
	idx := bytes.Index(data, separator)
	switch idx {
	case -1:
		return 0, false, nil
	case 0:
		return len(separator), true, nil
	}

	bytesParsed := idx

	currentLine := data[:idx]
	field_name, field_value, err := parseHeaderLine(currentLine)
	if err != nil {
		return 0, false, err
	}
	h[field_name] = field_value

	n, done, err := h.Parse(data[idx + len(separator):])
	bytesParsed += n + len(separator)
	if err != nil {
		return bytesParsed, false, err
	}

	return bytesParsed, done, nil
} 
