package headers

import (
	"bytes"
	"fmt"

	"github.com/oleshko-g/httpfromtcp/internal/stringio"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers, 0)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	const crlf = "\r\n"
	var bytesParsed int
	for crlfIndex := bytes.Index(data, []byte(crlf)); ; crlfIndex = bytes.Index(data, []byte(crlf)) {
		if crlfIndex > 0 { // probably a field line. try to parse
			rawFieldLine := data[:crlfIndex]
			colonIndex := bytes.Index(rawFieldLine, []byte(":"))
			if colonIndex > 0 {
				fieldName, err := validFieldName(rawFieldLine[:colonIndex])
				if err != nil {
					return 0, false, fmt.Errorf("400 Bad Request")
				}

				fieldValue := string(bytes.TrimSpace(rawFieldLine[colonIndex+1:]))
				if len(fieldValue) == 0 {
					bytesParsed = crlfIndex + 2 // emtry value. discard this rawFieldLine
					data = data[bytesParsed:]
					continue
				}

				if v, ok := h[fieldName]; ok {
					h[fieldName] = v + ", " + string(fieldValue)
					bytesParsed = crlfIndex + 2
					data = data[bytesParsed:]
					continue
				}

				h[fieldName] = string(fieldValue)
				bytesParsed = crlfIndex + 2
				data = data[bytesParsed:]
				continue
			}

			if colonIndex == -1 || colonIndex == 0 {
				return 0, false, fmt.Errorf("400 Bad Request")
			}
		}

		if crlfIndex == 0 { // end of field line section
			bytesParsed += 2
			done = true
			return bytesParsed, done, nil
		}

		if crlfIndex == -1 { // need more data
			return bytesParsed + 0, false, nil
		}
	}
}

func validFieldName(data []byte) (string, error) {
	fieldName := string(bytes.ToLower(data))
	for i, v := range fieldName {
		if stringio.IsWhiteSpace(v) {
			return "", fmt.Errorf("white space in the field name at index [%d]", i)
		}

		if !stringio.IsDigit(v) &&
			!stringio.IsLowerCaseLetter(v) &&
			!stringio.IsValidSpecialCharacter(v) {
			return "", fmt.Errorf("invalid character %v at index [%d]", v, i)
		}
	}

	return fieldName, nil
}
