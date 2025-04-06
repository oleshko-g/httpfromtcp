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
	crlfIndex := bytes.Index(data, []byte(crlf))
	if crlfIndex > 0 { // probably a field line. try to parse
		rawFieldLine := data[:crlfIndex]
		colonIndex := bytes.Index(rawFieldLine, []byte(":"))
		if colonIndex > 0 {
			fieldName := string(rawFieldLine[:colonIndex])
			if stringio.ContainsWhiteSpace(fieldName) {
				return 0, false, fmt.Errorf("400 Bad Request")
			}

			fieldValue := bytes.Fields(rawFieldLine[colonIndex+1:])
			if len(fieldValue) == 1 {
				h[fieldName] = string(fieldValue[0])
				return crlfIndex + 2, false, nil
			}

			if len(fieldValue) == 0 || len(fieldValue) > 1 {
				return crlfIndex + 2, false, nil // discard this rawFieldLine
			}
		}

		if colonIndex == -1 || colonIndex == 0 {
			return 0, false, fmt.Errorf("400 Bad Request")
		}
	}

	if crlfIndex == 0 { // end of field line section
		return 2, true, nil
	}

	if crlfIndex == -1 { // need more data
		return 0, false, nil
	}

	return 0, false, nil
}
