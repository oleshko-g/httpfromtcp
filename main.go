package main

import (
	"fmt"
	"io"
	"os"
)

const (
	messages = "./messages.txt"
)

func main() {
	file, err := os.Open(messages)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer file.Close()

	buf := make([]byte, 8)

	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			file.Close()
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "read: %s\n", buf[0:n])
	}
}
