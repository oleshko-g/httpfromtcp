package main

import (
	"fmt"
	"net"
	"os"

	"github.com/oleshko-g/httpfromtcp/internal/request"
)

const (
	network = "tcp"
	port    = ":42069"
)

func main() {
	address := fmt.Sprintf("127.0.0.1%s", port)
	listener, err := net.Listen(network, address)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Server listening on %s\n", address)
	defer func() {
		err := listener.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		fmt.Println("A connection has been accepted")
		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		req.PrintRequestLine()
		req.PrintHeaders()
		req.PrintBody()
		fmt.Println("A connection has been closed")
	}

}
