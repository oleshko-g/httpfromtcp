package main

import (
	"fmt"
	"io"
	"os"
)

const (
	messages = "./messages.txt"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	buf := make([]byte, 8)
	go func() {
		for line := make([]byte, 0); ; {
			n, err := f.Read(buf)
			if err == io.EOF {
				if len(line) != 0 {
					lines <- string(line)
				}
				f.Close()
				close(lines)
				break
			}

			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				f.Close()
				os.Exit(1)
			}

			for _, r := range buf[:n] {
				if r == '\n' {
					lines <- string(line)
					line = line[:0]
					continue
				}
				line = append(line, r)
			}
		}
	}()
	return lines
}

func main() {
	file, err := os.Open(messages)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	lines := getLinesChannel(file)

	for v := range lines {
		fmt.Fprintf(os.Stdout, "read: %s\n", v)
	}

}
