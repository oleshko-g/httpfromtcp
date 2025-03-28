package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	network = "udp"
	host    = "localhost"
	port    = "42069"
	address = host + ":" + port
)

func repl(r *bufio.Reader, c *net.UDPConn) error {
	for {
		fmt.Print(">")
		string, err := r.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		n, err := c.Write([]byte(string))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Printf("Wrote [%d] bytes\n", n)
	}
}

func main() {
	UDPAddr, err := net.ResolveUDPAddr(network, address)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Printf("Resolved adress [%s] to [%T] [%v]\n", address, UDPAddr, UDPAddr)

	conn, err := net.DialUDP(network, nil, UDPAddr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Printf("Established a [%s] connection [%#v] on [%v] address\n", network, conn, UDPAddr)
	defer conn.Close()

	err = repl(bufio.NewReader(os.Stdin), conn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
