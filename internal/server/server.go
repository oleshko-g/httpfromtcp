package server

import (
	"fmt"
	"net"
	"os"
	"sync/atomic"

	_ "github.com/oleshko-g/httpfromtcp/internal/http"
	"github.com/oleshko-g/httpfromtcp/internal/response"
)

const (
	network = "tcp"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	var server Server
	address := fmt.Sprintf("127.0.0.1:%d", port)
	var err error
	server.listener, err = net.Listen(network, address)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	go server.listen()

	return &server, nil
}

func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return err
	}
	s.closed.Store(true)
	return nil
}

func (s *Server) listen() {
	for {
		if s.closed.Load() {
			break
		}

		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		go s.handle(conn)

	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	headers := response.GetDefaultHeaders(0)
	response.WriteStatusLine(conn, response.StatusCodeOK())
	err := response.WriteHeaders(conn, headers)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
