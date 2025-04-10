package request

import "fmt"

func (r Request) PrintRequestLine() {
	fmt.Println("Request line:")
	fmt.Printf("- Method: %s\n- Target: %s\n- Version: %s\n", r.RequestLine.Method, r.RequestLine.RequestTarget, r.RequestLine.HttpVersion)
}

func (r Request) PrintHeaders() {
	fmt.Println("Headers:")
	for header, value := range r.Headers {
		fmt.Printf("- %s: %s\n", header, value)
	}
}
