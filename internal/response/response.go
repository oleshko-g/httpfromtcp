package response

import (
	"io"
	"strconv"

	"github.com/oleshko-g/httpfromtcp/internal/headers"
	"github.com/oleshko-g/httpfromtcp/internal/http"
)

var statusCodes = map[StatusCode]string{
	StatusCodeOK():                  "OK",
	StatusCodeBadRequest():          "Bad Request",
	StatusCodeInternalServerError(): "Internal Server Error",
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
	headers := headers.Headers{
		"Content-Length": strconv.Itoa(contentLength),
		"Connection":     "close",
		"Content-Type":   "text/plain",
	}
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
