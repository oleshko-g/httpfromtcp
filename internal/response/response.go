package response

import (
	"io"
	"strconv"

	"github.com/oleshko-g/httpfromtcp/internal/headers"
	_ "github.com/oleshko-g/httpfromtcp/internal/headers"
	"github.com/oleshko-g/httpfromtcp/internal/http"
	_ "github.com/oleshko-g/httpfromtcp/internal/server"
)

var statusCodes = map[StatusCode]string{
	statusCodeOK():                  "OK",
	statusCodeBadRequest():          "Bad Request",
	statusCodeInternalServerError(): "Internal Server Error",
}

type StatusCode [3]rune

func statusCodeOK() StatusCode {
	return StatusCode{'2', '0', '0'}
}

func statusCodeBadRequest() StatusCode {
	return StatusCode{'4', '0', '0'}
}

func statusCodeInternalServerError() StatusCode {
	return StatusCode{'5', '0', '0'}
}

type StatusLine string

func newStatusLine(version string, s StatusCode) string {
	reasonPhrase := statusCodes[s]
	return http.GetHttpVersion(version) + " " + string(s[:]) + " " + reasonPhrase + "\r\n"
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
		"Content-Type":   "test/plain",
	}
	return headers
}
