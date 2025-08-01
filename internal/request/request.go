package request

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

type ParseState int

const (
	StateInitialized ParseState = iota
	StateDone
)

type Request struct {
	RequestLine *RequestLine
	State ParseState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
	reqLine, count, err := parseRequestLine(data)
	if count != 0 && err == nil {
		r.RequestLine = reqLine
		r.State = StateDone
	} else if err != nil {
		r.State = StateDone
	}
	return count, err
}

func splitRequestLine(line string) (string, string, string, error) {
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 3 {
		return "", "", "", fmt.Errorf(("request line ill-formed: missing part"))
	}
	return parts[0], parts[1], parts[2], nil
}

func parseRequestLine(head []byte) (*RequestLine, int, error) {
	fmt.Println(head)
	line, _, found := strings.Cut(string(head), "\r\n")
	if !found {
		return nil, 0, nil
	}


	method, path, httpVer, err := splitRequestLine(line)
	if err != nil {
		return nil, 0, err
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
		return nil, 0, fmt.Errorf("request method contains invalid characters")
	}

	// Check http version
	ver, found := strings.CutPrefix(httpVer, "HTTP/")
	if !found || ver != "1.1" {
		return nil, 0, fmt.Errorf("request HTTP version not supported")
	}

	return &RequestLine{Method: method, RequestTarget: path, HttpVersion: ver}, len(line) + 2, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := Request {RequestLine: nil, State: StateInitialized}
	data := make([]byte, 0)


	buf := make([]byte, 8)
	for req.State != StateDone {
		read, err := reader.Read(buf)
		if err != nil {
			return &req, fmt.Errorf("error while reading from stream: %w", err)
		}
		data = slices.Concat(data, buf[:read])
		parsed, err := req.parse(data)
		if err != nil {
			return &req, fmt.Errorf("error while parsing request: %w", err)
		}
		if parsed != 0 {
			data = data[parsed:]
		}
	}

	for {
		_, err := reader.Read(buf)
		if err == io.EOF {
			return &req, nil
		} else if err != nil {
			return &req, fmt.Errorf("error while reading from stream: %w", err)
		}
	}
}