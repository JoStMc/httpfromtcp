package headers

import (
	"bytes"
	"errors"
	"slices"
	"strings"
)

type Headers map[string]string


func NewHeaders() Headers {
	return make(Headers)
} 

var separator = []byte("\r\n")
var validSpecialChars = []rune{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'} 


func isToken(str string) bool {
	for _, ch := range str {
	    if ch < 'A' || ch > 'Z' && ch < 'a' || ch > 'z' {
	        if ch < '0' || ch > '9' {
				if !slices.Contains(validSpecialChars, ch) {
					return false
				} 
	        } 
	    } 
	} 
	return true
} 

func parseHeaderLine(line []byte) (string, string, error) {
	parts := bytes.Fields(line)
	if len(parts) != 2 {
		return "", "", errors.New("invalid header line")
	} 
	fieldName := string(parts[0])
	fieldValue := string(parts[1])
	if fieldName[len(fieldName)-1] != ':' {
		return "", "", errors.New("invalid header line")
	} 
	fieldName = fieldName[:len(fieldName)-1]

	if !isToken(fieldName) {
		return "", "", errors.New("field name contains invalid character")
	} 

	return fieldName, fieldValue, nil
} 


func (h Headers) Get(name string) string {
    return h[strings.ToLower(name)]
} 

func (h Headers) Replace(name, value string) {
	h[strings.ToLower(name)] = value
} 

func (h Headers) Set(name, value string) {
	name = strings.ToLower(name)
	v, ok := h[name]
	if ok {
		value = v + "," + value
	} 
	h[name] = value
} 

// UGLY!!!
func (h Headers) ForEach(cb func(n, v string)) {
	for n, v := range h {
		cb(n, v)
	} 
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
	fieldName, fieldValue, err := parseHeaderLine(currentLine)
	if err != nil {
		return 0, false, err
	}

	h.Set(fieldName, fieldValue)

	n, done, err := h.Parse(data[idx + len(separator):])
	bytesParsed += n + len(separator)
	if err != nil {
		return bytesParsed, false, err
	}

	return bytesParsed, done, nil
} 
