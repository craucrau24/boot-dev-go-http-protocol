package response

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/craucrau24/boot-dev-go-http-protocol/internal/headers"
)

type writerStep int
const (WriterStatusLineStep writerStep = iota
	WriterHeadersStep
	WriterBodyStep
	WriterEndStep
)

type Writer struct {
	writer io.Writer
	step writerStep
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.step == WriterStatusLineStep {
		w.step = WriterHeadersStep
		return WriteStatusLine(w.writer, statusCode)
	} else {
		return fmt.Errorf("attempt to override status line (second call)")
	}
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.step == WriterHeadersStep {
		w.step = WriterBodyStep
		return WriteHeaders(w.writer, headers)
	} else {
		return fmt.Errorf("headers should follow immediately status line")
	}
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.step == WriterBodyStep {
		w.step = WriterBodyStep
		return w.writer.Write(p)
	} else {
		return 0, fmt.Errorf("body should follow immediately headers")
	}
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	size := strings.ToUpper(strconv.FormatUint(uint64(len(p)), 16))
	return w.writer.Write([]byte(fmt.Sprintf("%s\r\n%s\r\n", size, p)))
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return w.WriteChunkedBody(nil)
}


func NewWriter(writer io.Writer) Writer {
	return Writer {writer: writer, step: WriterStatusLineStep}
}