package request

import (
	"fmt"
	"io"
	"strings"
)

var versionsSupported map[string]struct{} = map[string]struct{}{
	"HTTP/1.1": {},
}

var methodsSupported map[string]struct{} = map[string]struct{}{
	"GET": {},
}

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	Method        string
	RequestTarget string
}

func upperCaseLetters(s string) bool {
	for _, rune := range s {
		// 'A' is a literal of 0x41 number in ASCII encoding
		// 'Z' is literal of 0x5A number in ASCII encoding
		if rune < 'A' || rune > 'Z' {
			return false
		}
	}
	return true
}

func containsWhiteSpace(s string) bool {
	return len(strings.Fields(s)) > 1
}

func methodSupported(s string) bool {
	_, ok := methodsSupported[s]
	return ok
}

func versionSupported(s string) bool {
	_, ok := versionsSupported[s]
	return ok
}

func RequestFromReader(r io.Reader) (*Request, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	requestLines := strings.Split(string(data), "\r\n")

	parts := strings.Split(string(requestLines[0]), " ")
	if len(parts) < 3 {
		return nil, fmt.Errorf("400 Bad Request")
	}

	method, target, version := parts[0], parts[1], parts[2]

	if !upperCaseLetters(method) {
		return nil, fmt.Errorf("400 Bad Request")
	}

	if containsWhiteSpace(target) {
		return nil, fmt.Errorf("400 Bad Request")
	}

	if !versionSupported(version) {
		return nil, fmt.Errorf("505 HTTP Version Not Supported")
	}

	if !methodSupported(method) {
		return nil, fmt.Errorf("501 Not Implemented")
	}

	return &Request{
		RequestLine: RequestLine{
			HttpVersion:   strings.Split(version, "/")[1],
			Method:        method,
			RequestTarget: target,
		},
	}, nil
}
