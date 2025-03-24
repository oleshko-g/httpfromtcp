package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

const (
	network = "tcp"
)

func getLinesChannel(r io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer close(lines)
		buf := make([]byte, 8)
		for line := make([]byte, 0); ; {
			n, err := r.Read(buf)
			if n > 0 {
				for _, r := range buf[:n] {
					if r == '\n' {
						lines <- string(line)
						line = line[:0]
						continue
					}
					line = append(line, r)
				}
			}

			if err == io.EOF {
				if len(line) != 0 {
					lines <- string(line)
				}
				err := r.Close()
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				break
			}

			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				err := r.Close()
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				os.Exit(1)
			}
		}
	}()

	return lines
}

func main() {
	listener, err := net.Listen(network, "127.0.0.1:42069")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer func() {
		err := listener.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	for {
		conn, err := listener.Accept()
		fmt.Println("A connection has been accepted")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		lines := getLinesChannel(conn)
		for v := range lines {
			fmt.Fprintf(os.Stdout, "%s\n", v)
		}
		fmt.Println("A connection has been closed")
	}

}
