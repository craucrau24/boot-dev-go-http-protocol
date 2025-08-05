package request

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/craucrau24/boot-dev-go-http-protocol/internal/headers"
)

type ParseState int

const (
	StateInitialized ParseState = iota
	StateParsingHeaders
	StateDone
)

type Request struct {
	RequestLine *RequestLine
	Headers headers.Headers
	State ParseState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
	var count int
	var err error

	switch r.State {
	case StateInitialized: {
		var reqLine *RequestLine
		reqLine, count, err = parseRequestLine(data)
		if count != 0 && err == nil {
			r.RequestLine = reqLine
			r.State = StateParsingHeaders
		} else if err != nil {
			r.State = StateParsingHeaders
		}
	}

	case StateParsingHeaders: {
		var done bool
		count, done, err = r.Headers.Parse(data)
		// fmt.Printf("%d, %v, %v\n", count, done, err)
		if done {
			r.State = StateDone
		}
	}
	}
	// fmt.Printf("%d, %v\n", count, err)
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
	req := Request {RequestLine: nil, State: StateInitialized, Headers: headers.NewHeaders()}
	data := make([]byte, 0)
	var err error


	buf := make([]byte, 8)
	for req.State != StateDone && err == nil {
		var read int
		read, err = reader.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("error while reading from stream: %w", err)
		}
		data = slices.Concat(data, buf[:read])
		for {
			read, err = req.parse(data)
			// fmt.Printf("read: %v, err: %v\n", read, err)
			if err != nil {
				return nil, fmt.Errorf("error while parsing request: %w", err)
			}
			if read != 0 {
				data = data[read:]
			} else {
				break
			}
		}
	}

	if req.State == StateDone {
		return &req, nil
	} else {
		return nil, fmt.Errorf("EOF reached: truncated request")
	}

	// for {
	// 	_, err := reader.Read(buf)
	// 	if err == io.EOF {
	// 		return &req, nil
	// 	} else if err != nil {
	// 		return nil, fmt.Errorf("error while reading from stream: %w", err)
	// 	}
	// }
}