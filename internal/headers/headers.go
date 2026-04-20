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
	parts := bytes.Fields(currentLine)
	if len(parts) != 2 {
		return 0, false, errors.New("invalid header line")
	} 
	field_name := string(parts[0])
	field_value := string(parts[1])
	if field_name[len(field_name)-1] != ':' {
		return 0, false, errors.New("invalid header line")
	} 
	h[strings.TrimRight(field_name, ":")] = field_value

	n, done, err := h.Parse(data[idx + len(separator):])
	bytesParsed += n
	if err != nil {
		return bytesParsed, false, err
	}

	return bytesParsed, done, nil
} 
