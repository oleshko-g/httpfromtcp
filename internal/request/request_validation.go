package request

import (
	"github.com/oleshko-g/httpfromtcp/internal/stringio"
)

const (
	httpName = "HTTP"
)

func validHTTPVersion(s string) bool {

	if len(s) != 8 { // lenth of 'HTTP/DIGIT.DIGIT" in bytes
		return false
	}

	if s[:4] != httpName {
		return false
	}

	if s[4] != '/' {
		return false
	}

	if !stringio.IsDigit(rune(s[5])) {
		return false
	}

	if s[6] != '.' {
		return false
	}

	if !stringio.IsDigit(rune(s[7])) {
		return false
	}

	return true
}

func validHTTPMethod(s string) bool {
	return stringio.UpperCaseLetters(s)
}

func validHTTPTarget(s string) bool {
	return !stringio.ContainsWhiteSpace(s)
}
