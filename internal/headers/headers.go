package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	header, _, found := bytes.Cut(data, []byte{'\r', '\n'})
	count := len(header) + 2
	if !found {
		return 0, false, nil
	}

	if len(header) == 0{
		return 2, true, nil
	}

	header = bytes.Trim(header, " ")
	name, value, _ := bytes.Cut(header, []byte{':'})
	if bytes.IndexFunc(name, invalidRune) >= 0 {
		return 0, false, fmt.Errorf("header name contains invalid character")
	}
	h[strings.ToLower(string(name))] = string(bytes.TrimLeft(value, " "))
	return count, false, nil
}

func invalidRune(r rune) bool {
  allowedChar := []rune{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}
	switch {
	case r >= 'A' && r <= 'Z':
		return false
	case r >= 'a' && r <= 'z':
		return false
	case r >= '0' && r <= '9':
		return false
	}
	return !slices.Contains(allowedChar, r)
}

func NewHeaders() Headers {
	return make(map[string]string)
}