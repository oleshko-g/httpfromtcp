package request

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/oleshko-g/httpfromtcp/internal/headers"
	"github.com/oleshko-g/httpfromtcp/internal/http"
)

type RequestState string

func RequestStateInitialized() RequestState {
	return RequestState("initialized")
}

func RequestStateParsingHeaders() RequestState {
	return RequestState("parsing headers")
}

func RequestStateParsingBody() RequestState {
	return RequestState("parsing body")
}

func RequestStateDone() RequestState {
	return RequestState("done")
}

type Request struct {
	RequestLine   RequestLine
	Headers       headers.Headers
	Body          []byte
	state         RequestState
	contentLength int
}

func (r *Request) getContentLength() (int, bool, error) {
	value, ok := r.Headers.Get("content-length")
	if ok {
		contentLength, err := strconv.Atoi(value)
		if err != nil {
			return 0, false, err
		}
		return contentLength, ok, nil
	}

	return 0, false, nil
}

func (r *Request) parse(data []byte) (int, error) {
	var bytesParsed int
	var err error
	switch r.state {
	case RequestStateInitialized():
		bytesParsed, r.RequestLine, err = parseRequestLine(data)
		if bytesParsed > 0 {
			defer func() { r.state = RequestStateParsingHeaders() }()
			return bytesParsed, nil
		}
	case RequestStateParsingHeaders():
		bytesParsed, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			contentLength, ok, err := r.getContentLength()
			if err != nil {
				return 0, err
			}
			if !ok {
				r.state = RequestStateDone()
				return 0, nil
			}
			r.contentLength = contentLength
			r.state = RequestStateParsingBody()
		}
		return bytesParsed, nil
	case RequestStateParsingBody():

		bytesRemaining := r.contentLength - len(r.Body) // also handles zero content length

		bytesToAppend := min(bytesRemaining, len(data))

		r.Body = append(r.Body, data[:bytesToAppend]...) // also handles zero content length – appends up to zero bytes
		bytesParsed = bytesToAppend

		if len(r.Body) == r.contentLength {
			r.state = RequestStateDone()
		}
		return bytesParsed, nil
	}
	return 0, err
}

type RequestLine struct {
	HttpVersion   string
	Method        string
	RequestTarget string
}

func RequestFromReader(r io.Reader) (*Request, error) {
	request := Request{
		RequestLine: RequestLine{},
		Headers:     make(headers.Headers),
		state:       RequestStateInitialized(),
	}

	var bytesReadTo int
	for buf := make([]byte, 8); request.state != RequestStateDone(); {
		if bytesReadTo == len(buf) {
			buffer := make([]byte, len(buf)*2)
			copy(buffer, buf)
			buf = buffer
		}
		bytesRead, errRead := r.Read(buf[bytesReadTo:])
		if errRead != nil && errRead != io.EOF {
			return &Request{}, errRead
		}
		bytesReadTo += bytesRead

		bytesParsed, errParse := request.parse(buf[:bytesReadTo])
		if errParse != nil {
			return &Request{}, errParse
		}
		if bytesParsed > 0 {
			copy(buf, buf[bytesParsed:bytesReadTo])
			bytesReadTo -= bytesParsed
		}

		if errRead == io.EOF {
			if request.state != RequestStateDone() {
				return nil, fmt.Errorf("unexpected EOF before complete request")
			}
			break
		}
	}

	return &request, nil
}

func parseRequestLine(data []byte) (int, RequestLine, error) {
	const crlf = "\r\n"
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, RequestLine{}, nil
	}
	requestString := string(data[:idx])

	parts := strings.Split(requestString, " ")
	if len(parts) != 3 {
		return 0, RequestLine{}, fmt.Errorf("400 Bad Request 3")
	}

	if !http.ValidHTTPMethod(parts[0]) {
		return 0, RequestLine{}, fmt.Errorf("400 Bad Request 4")
	}
	method := parts[0]

	if !http.ValidHTTPTarget(parts[1]) {
		return 0, RequestLine{}, fmt.Errorf("400 Bad Request 5")
	}
	target := parts[1]

	if !http.ValidHTTPVersion(parts[2]) {
		return 0, RequestLine{}, fmt.Errorf("400 Bad Request 6")
	}
	version := strings.Split(parts[2], "/")[1]

	if !VersionSupported(version) {
		return 0, RequestLine{}, fmt.Errorf("505 HTTP Version Not Supported")
	}

	if !MethodSupported(method) {
		return 0, RequestLine{}, fmt.Errorf("501 Not Implemented")
	}
	return idx + 2, RequestLine{
		HttpVersion:   version,
		Method:        method,
		RequestTarget: target,
	}, nil
}
