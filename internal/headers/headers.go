package headers

import (
	"bytes"
	"fmt"
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
	if bytes.ContainsAny(name, " ") {
		return 0, false, fmt.Errorf("header name shall not contain spaces")
	}
	h[string(name)] = string(bytes.TrimLeft(value, " "))
	return count, false, nil
}

func NewHeaders() Headers {
	return make(map[string]string)
}