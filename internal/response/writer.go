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
		w.step = WriterEndStep
		return w.writer.Write(p)
	} else {
		return 0, fmt.Errorf("body should follow immediately headers")
	}
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.step == WriterBodyStep {
		size := strings.ToUpper(strconv.FormatUint(uint64(len(p)), 16))
		return w.writer.Write([]byte(fmt.Sprintf("%s\r\n%s\r\n", size, p)))
	} else {
		return 0, fmt.Errorf("body should follow immediately headers")
	}
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
		n, err := w.WriteChunkedBody(nil)
		w.step = WriterEndStep
		return n, err
}


func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.step == WriterBodyStep {
		w.step = WriterEndStep
		w.writer.Write([]byte("0\r\n"))
		return WriteHeaders(w.writer, h)
	} else {
		return fmt.Errorf("trailers should follow immediately body")
	}
}

func NewWriter(writer io.Writer) Writer {
	return Writer {writer: writer, step: WriterStatusLineStep}
}