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
	for line := make([]byte, 0); ; {
		n, err := file.Read(buf)
		if err == io.EOF {
			if len(line) != 0 {
				fmt.Fprintf(os.Stdout, "read: %s\n", string(line))
			}
			break
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			file.Close()
			os.Exit(1)
		}

		for _, r := range buf[:n] {
			if r == '\n' {
				fmt.Fprintf(os.Stdout, "read: %s\n", string(line))
				line = line[:0]
				continue
			}
			line = append(line, r)
		}
	}
}
