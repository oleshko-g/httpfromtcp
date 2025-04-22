package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/oleshko-g/httpfromtcp/internal/headers"
	"github.com/oleshko-g/httpfromtcp/internal/http"
)

var statusCodes = map[StatusCode]string{
	StatusCodeOK():                  "OK",
	StatusCodeBadRequest():          "Bad Request",
	StatusCodeNotFound():            "Not Found",
	StatusCodeInternalServerError(): "Internal Server Error",
}

type writerState string

func writerStateInitialized() writerState {
	return "Initialized"
}
func writerStateStatusLineWritten() writerState {
	return "Status line written"
}

func writerStateHeadersWritten() writerState {
	return "Headers written"
}

func writerStateBodyWritingStarted() writerState {
	return "Body writing have started"
}

func writerStateBodyWritten() writerState {
	return "Body written"
}

func writerStateDone() writerState {
	return "Done"
}

type Writer struct {
	conn  io.Writer
	state writerState
}

func NewWriter(conn io.Writer) *Writer {
	return &Writer{
		state: writerStateInitialized(),
		conn:  conn,
	}
}

func (w *Writer) WriteStatusLine(sc StatusCode) error {
	if w.state != writerStateInitialized() {
		return fmt.Errorf("trying to WriteStatusLine() int Initialized Writer state")
	}
	buf := []byte(newStatusLine("1.1", sc))
	_, err := w.conn.Write(buf)
	if err == nil {
		w.state = writerStateStatusLineWritten()
	}
	return err
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.state != writerStateStatusLineWritten() {
		return fmt.Errorf("trying to WriteHeaders() not in StatusLineWritten state")
	}
	buf := headersToBuf(h)
	_, err := w.conn.Write(buf)
	if err == nil {
		w.state = writerStateHeadersWritten()
		if v, _ := h.Get("transfer-encoding"); v == "chunked" {
			return nil
		}
		if _, ok := h.Get("content-length"); !ok {
			w.state = writerStateDone()
			return nil
		}
	}
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != writerStateHeadersWritten() {
		return 0, fmt.Errorf("trying to WriteBody() not in HeadersWritten state")
	}
	n, err := w.conn.Write(p)
	if err == nil {
		w.state = writerStateBodyWritten()
		w.state = writerStateDone()
	}
	return n, err
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.state != writerStateHeadersWritten() && w.state != writerStateBodyWritingStarted() {
		return 0, fmt.Errorf("trying to WriteChunkedBody() in the states that are not HeadersWritten, BodyWritingStarted")
	}
	p = wrapChunk(p)
	n, err := w.conn.Write(p)
	if err == nil {
		w.state = writerStateBodyWritingStarted()
		return n, nil
	}
	return n, err
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.state != writerStateHeadersWritten() && w.state != writerStateBodyWritingStarted() {
		return 0, fmt.Errorf("trying to WriteChunkedBodyDone() in the states that are not HeadersWritten, BodyWritingStarted")
	}
	var emptySlice []byte
	lastChunk := wrapChunk(emptySlice)
	n, err := w.conn.Write(lastChunk)
	if err == nil {
		w.state = writerStateBodyWritten()
		w.state = writerStateDone()
	}
	return n, err
}

func wrapChunk(c []byte) []byte {
	hexChunkLength := fmt.Sprintf("%X", len(c))
	bytePrefix := append([]byte(hexChunkLength), http.ByteCRLF...)
	bytePostfix := http.ByteCRLF
	prefixedChunk := append(bytePrefix, c...)
	wrappedChunk := append(prefixedChunk, bytePostfix...)

	return wrappedChunk
}

type StatusCode [3]rune

func (s *StatusCode) String() string {
	return string(s[:])
}

func StatusCodeOK() StatusCode {
	return StatusCode{'2', '0', '0'}
}

func StatusCodeBadRequest() StatusCode {
	return StatusCode{'4', '0', '0'}
}

func StatusCodeNotFound() StatusCode {
	return StatusCode{'4', '0', '4'}
}

func StatusCodeInternalServerError() StatusCode {
	return StatusCode{'5', '0', '0'}
}

func newStatusLine(version string, s StatusCode) string {
	reasonPhrase := statusCodes[s]
	return http.GetHttpVersion(version) + " " +
		s.String() + " " +
		reasonPhrase + "\r\n"
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	buf := []byte(newStatusLine("1.1", statusCode))
	_, err := w.Write(buf)
	return err
}

func GetDefaultHeaders(contentLength int) headers.Headers {
	headers := headers.Headers{}
	headers.Set("Content-Length", strconv.Itoa(contentLength))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	buf := headersToBuf(headers)
	_, err := w.Write(buf)
	return err
}

func headersToBuf(headers headers.Headers) []byte {
	var buf []byte
	for header, value := range headers {
		buf = append(buf, []byte(header+": ")...)
		buf = append(buf, []byte(value)...)
		buf = append(buf, '\r', '\n')
	}
	buf = append(buf, '\r', '\n')
	return buf
}

func WriteBody(w io.Writer, s string) {
	buf := []byte(s)
	buf = append(buf, '\r', '\n')
	w.Write(buf)
}
