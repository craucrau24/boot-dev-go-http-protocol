package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// fmt.Printf("Headers: %s\n", data)
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
	val := string(bytes.TrimLeft(value, " "))
	key := string(name)
	old, ok := h.Get(key)
	if ok {
		h.Set(key, old + ", " + val)
	} else {
		h.Set(key, val)
	}
	// fmt.Printf("count: %d\n", count)
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

func (h Headers) Get(name string) (string, bool) {
	key := strings.ToLower(name)
	val, ok := h[key]
	return val, ok
}

func (h Headers) Set(name, value string) {
	key := strings.ToLower(name)
	h[key] = value
}

func (h Headers) Append(name, value string) {
	old, ok := h.Get(name)
	if ok {
		value = fmt.Sprintf("%s, %s", old, value)
	}
	h.Set(name, value)
}


func (h Headers) Unset(name string) {
	key := strings.ToLower(name)
	delete(h, key)
}

func NewHeaders() Headers {
	return make(map[string]string)
}