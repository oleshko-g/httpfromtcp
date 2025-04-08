package stringio

func UpperCaseLetters(s string) bool {
	for _, rune := range s {
		// 'A' is a literal of 0x41 number in ASCII encoding
		// 'Z' is literal of 0x5A number in ASCII encoding
		if rune < 'A' || rune > 'Z' {
			return false
		}
	}
	return true
}

func IsLowerCaseLetter(r rune) bool {
	return r >= 'a' && r <= 'z'
}
func IsValidSpecialCharacter(r rune) bool {
	isValid := false
	switch r {
	case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
		isValid = true
	}
	return isValid
}

func IsDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func ContainsWhiteSpace(s string) bool {
	for _, v := range s {
		if IsWhiteSpace(v) {
			return true
		}
	}
	return false
}

func IsWhiteSpace(r rune) bool {
	switch {
	case r == 0x0020: // Space
		return true
	case r == 0x00A0: // Non-breaking Space (NBSP)
		return true
	case r == '\r': // Carriage return
		return true
	case r == '\n': // Line feed
		return true
	case r == '\t': // Tab
		return true
	case r == '\f': // Form feed
		return true
	case r == '\v': // Vertical Tab
		return true
	}

	return false
}
