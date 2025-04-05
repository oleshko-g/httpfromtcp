package headers

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers, 0)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	return 0, false, nil
}
