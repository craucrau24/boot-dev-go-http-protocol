package response

import (
	"fmt"
	"io"
	"strings"

	"github.com/craucrau24/boot-dev-go-http-protocol/internal/headers"
)

type StatusCode int
const StatusOk StatusCode = 200
const StatusBadRequest StatusCode = 400
const StatusInternalServerError StatusCode = 500

var reasonMap = map[StatusCode]string{StatusOk: "OK", StatusBadRequest: "Bad Request", StatusInternalServerError: "Internal Server Error"}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	reason := reasonMap[statusCode]
  statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reason)
	_, err := w.Write([]byte(statusLine))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	headers.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	return headers
}

func WriteHeaders(w io.Writer, h headers.Headers) error {
	for name, value := range h {
		var build strings.Builder
		build.WriteString(name)
		build.WriteString(": ")
		build.WriteString(value)
		build.WriteString("\r\n")
		_, err := w.Write([]byte(build.String()))
		if err != nil {
			return fmt.Errorf("error while writing header: %w", err)
		}
	}
	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return fmt.Errorf("error while writing header: %w", err)
	}
	return nil
}
