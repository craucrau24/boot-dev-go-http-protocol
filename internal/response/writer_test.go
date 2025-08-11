package response

import (
	"bytes"
	"testing"

	"github.com/craucrau24/boot-dev-go-http-protocol/internal/headers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrated(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	err := w.WriteStatusLine(StatusOk)
	require.NoError(t, err)
	assert.Equal(t, "HTTP/1.1 200 OK\r\n", buf.String())

	err = w.WriteHeaders(GetDefaultHeaders(12))
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "content-type: text/plain\r\n")
	assert.Contains(t, buf.String(), "content-length: 12\r\n")
	assert.Equal(t, "\r\n\r\n", buf.String()[buf.Len() - 4:])
	buf.Reset()

	n, err := w.WriteChunkedBody([]byte("Foo is awesome!"))
	require.NoError(t, err)
	assert.Equal(t, 20, n)
	n, err = w.WriteChunkedBody([]byte("Bar is way better!"))
	require.NoError(t, err)
	assert.Equal(t, 24, n)
	n, err = w.WriteChunkedBody([]byte("Baz is epic loser!!"))
	require.NoError(t, err)
	assert.Equal(t, 25, n)
	assert.Equal(t, "F\r\nFoo is awesome!\r\n12\r\nBar is way better!\r\n13\r\nBaz is epic loser!!\r\n", buf.String())
	buf.Reset()

	head := headers.NewHeaders()
	head.Set("baz", "FooBar")
	err = w.WriteTrailers(head)
	require.NoError(t, err)
	assert.Equal(t, "0\r\nbaz: FooBar\r\n\r\n", buf.String())
}

