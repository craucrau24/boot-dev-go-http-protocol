package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func splitRequestLine(line string) (string, string, string, error) {
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 3 {
		return "", "", "", fmt.Errorf(("request line ill-formed: missing part"))
	}
	return parts[0], parts[1], parts[2], nil
}

func parseRequestLine(head []byte) (*RequestLine, error) {
	line, _, found := strings.Cut(string(head), "\r\n")
	if !found {
		return nil, fmt.Errorf("request line ill-formed: missing CRLF")
	}


	method, path, httpVer, err := splitRequestLine(line)
	if err != nil {
		return nil, err
	}

	// Check method
	allUpper := true
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			allUpper = false
			break
		}
	}
	if !allUpper {
		return nil, fmt.Errorf("request method contains invalid characters")
	}

	// Check http version
	ver, found := strings.CutPrefix(httpVer, "HTTP/")
	if !found || ver != "1.1" {
		return nil, fmt.Errorf("request HTTP version not supported")
	}

	return &RequestLine{Method: method, RequestTarget: path, HttpVersion: ver}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	reqLine, err := parseRequestLine(req)
	if err != nil {
		return nil, err
	}

	return &Request{RequestLine: *reqLine}, nil

}