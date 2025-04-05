package request

import (
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

func (rs *RequestState) RequestStateInitialized() RequestState {
	return RequestState("initialized")
}

func (rs *RequestState) RequestStateDone() RequestState {
	return RequestState("done")
}

type Request struct {
	RequestLine RequestLine
	state       RequestState
}

func (r *Request) parse(buf []byte) {

}

type RequestLine struct {
	HttpVersion   string
	Method        string
	RequestTarget string
}

func RequestFromReader(r io.Reader) (*Request, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	requestLines := strings.Split(string(data), "\r\n")

	requestLine, err := parseRequestLine([]byte(requestLines[0]))
	if err != nil {
		return &Request{}, err
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	parts := strings.Split(string(data), " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("400 Bad Request")
	}

	if !validHTTPMethod(parts[0]) {
		return nil, fmt.Errorf("400 Bad Request")
	}
	method := parts[0]

	if !validHTTPTarget(parts[1]) {
		return nil, fmt.Errorf("400 Bad Request")
	}
	target := parts[1]

	if !validHTTPVersion(parts[2]) {
		return nil, fmt.Errorf("400 Bad Request")
	}
	version := strings.Split(parts[2], "/")[1]

	if !versionSupported(version) {
		return nil, fmt.Errorf("505 HTTP Version Not Supported")
	}

	if !methodSupported(method) {
		return nil, fmt.Errorf("501 Not Implemented")
	}
	return &RequestLine{
		HttpVersion:   version,
		Method:        method,
		RequestTarget: target,
	}, nil
}
