package request

import (
	"strings"
)

const (
	httpName = "HTTP"
)

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

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

	if !isDigit(rune(s[5])) {
		return false
	}

	if s[6] != '.' {
		return false
	}

	if !isDigit(rune(s[7])) {
		return false
	}

	return true
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

func validHTTPMethod(s string) bool {
	return upperCaseLetters(s)
}

func containsWhiteSpace(s string) bool {
	return len(strings.Fields(s)) > 1
}

func validHTTPTarget(s string) bool {
	return !containsWhiteSpace(s)
}
