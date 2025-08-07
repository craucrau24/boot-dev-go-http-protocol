package response

import (
	"fmt"
	"strings"
	"testing"

	"github.com/craucrau24/boot-dev-go-http-protocol/internal/headers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func first(f string, ok bool) string {
	return f
}

func parseHeaders(content string, count int) (headers.Headers, error) {
	head := headers.NewHeaders()
	buf := []byte(content)
	for _ = range count {
		c, _, _ := head.Parse(buf)
		buf = buf[c:]
	}
	_, done, err := head.Parse(buf)
	if done && err == nil {
		return head, nil
	}
	return head, fmt.Errorf("unexpected EOF")
}

func TestWriteStatusLine(t *testing.T) {
	var builder strings.Builder
	err := WriteStatusLine(&builder, StatusOk)
	require.NoError(t, err)
	assert.Equal(t, "HTTP/1.1 200 OK\r\n", builder.String())

	builder.Reset()
	err = WriteStatusLine(&builder, StatusBadRequest)
	require.NoError(t, err)
	assert.Equal(t, "HTTP/1.1 400 Bad Request\r\n", builder.String())

	builder.Reset()
	err = WriteStatusLine(&builder, StatusInternalServerError)
	require.NoError(t, err)
	assert.Equal(t, "HTTP/1.1 500 Internal Server Error\r\n", builder.String())

	builder.Reset()
	err = WriteStatusLine(&builder, 302)
	require.NoError(t, err)
	assert.Equal(t, "HTTP/1.1 302 \r\n", builder.String())
}

func TestGetDefaultHeaders(t *testing.T) {
	headers := GetDefaultHeaders(0)
	assert.Equal(t, "close", first(headers.Get("Connection")))
	assert.Equal(t, "text/plain", first(headers.Get("Content-Type")))
	assert.Equal(t, "0", string(first(headers.Get("Content-Length"))))

	headers = GetDefaultHeaders(10)
	assert.Equal(t, "close", first(headers.Get("Connection")))
	assert.Equal(t, "text/plain", first(headers.Get("Content-Type")))
	assert.Equal(t, "10", string(first(headers.Get("Content-Length"))))

	headers = GetDefaultHeaders(42)
	assert.Equal(t, "close", first(headers.Get("Connection")))
	assert.Equal(t, "text/plain", first(headers.Get("Content-Type")))
	assert.Equal(t, "42", string(first(headers.Get("Content-Length"))))
}

func TestWriteHeaders(t *testing.T) {
	var builder strings.Builder
	head := GetDefaultHeaders(0)
	err := WriteHeaders(&builder, head)
	require.NoError(t, err)
	expected, err := parseHeaders("Connection: close\r\nContent-Type: text/plain\r\nContent-Length: 0\r\n\r\n", 3)
	require.NoError(t, err)
	actual, err := parseHeaders(builder.String(), 3)
	require.NoError(t, err)
	for _, key := range []string{"Connection", "Content-Type", "Content-Length"} {
		assert.Equal(t, first(expected.Get(key)), first(actual.Get(key)))
	}

	builder.Reset()
	head = GetDefaultHeaders(42)
	err = WriteHeaders(&builder, head)
	require.NoError(t, err)
	expected, err = parseHeaders("Connection: close\r\nContent-Type: text/plain\r\nContent-Length: 42\r\n\r\n", 3)
	require.NoError(t, err)
	actual, err = parseHeaders(builder.String(), 3)
	require.NoError(t, err)
	for _, key := range []string{"Connection", "Content-Type", "Content-Length"} {
		assert.Equal(t, first(expected.Get(key)), first(actual.Get(key)))
	}
}