package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

var versionsSupported = map[string]struct{}{
	"1.1": {},
}

func versionSupported(s string) bool {
	_, ok := versionsSupported[s]
	return ok
}

var methodsSupported = map[string]struct{}{
	"GET": {},
}

func methodSupported(s string) bool {
	_, ok := methodsSupported[s]
	return ok
}

type RequestState string

func RequestStateInitialized() RequestState {
	return RequestState("initialized")
}

func RequestStateDone() RequestState {
	return RequestState("done")
}

type Request struct {
	RequestLine RequestLine
	state       RequestState
}

func (r *Request) parse(data []byte) (int, error) {
	n, rl, err := parseRequestLine(data)
	if n > 0 {
		r.RequestLine = rl
		r.state = RequestStateDone()
	}
	return n, err
}

type RequestLine struct {
	HttpVersion   string
	Method        string
	RequestTarget string
}

func RequestFromReader(r io.Reader) (*Request, error) {
	request := Request{
		RequestLine: RequestLine{},
		state:       RequestStateInitialized(),
	}

	var bytesReadTo int
	var bytesParsedTo int
	for buf := make([]byte, 8); request.state != RequestStateDone(); {
		if bytesReadTo == len(buf) {
			buffer := make([]byte, len(buf)*2)
			copy(buffer, buf)
			buf = buffer
		}

		bytesRead, errRead := r.Read(buf[bytesReadTo:])
		if errRead != nil {
			return &Request{}, errRead
		}
		bytesReadTo += bytesRead

		bytesParsed, errParse := request.parse(buf[:bytesReadTo])
		if errParse != nil {
			return &Request{}, errParse
		}
		if bytesParsed > 0 {
			bytesParsedTo += bytesParsed
			copy(buf, buf[bytesParsedTo:bytesReadTo])
			bytesReadTo -= bytesParsedTo
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
		return 0, RequestLine{}, fmt.Errorf("400 Bad Request")
	}

	if !validHTTPMethod(parts[0]) {
		return 0, RequestLine{}, fmt.Errorf("400 Bad Request")
	}
	method := parts[0]

	if !validHTTPTarget(parts[1]) {
		return 0, RequestLine{}, fmt.Errorf("400 Bad Request")
	}
	target := parts[1]

	if !validHTTPVersion(parts[2]) {
		return 0, RequestLine{}, fmt.Errorf("400 Bad Request")
	}
	version := strings.Split(parts[2], "/")[1]

	if !versionSupported(version) {
		return 0, RequestLine{}, fmt.Errorf("505 HTTP Version Not Supported")
	}

	if !methodSupported(method) {
		return 0, RequestLine{}, fmt.Errorf("501 Not Implemented")
	}
	return idx + 2, RequestLine{
		HttpVersion:   version,
		Method:        method,
		RequestTarget: target,
	}, nil
}
