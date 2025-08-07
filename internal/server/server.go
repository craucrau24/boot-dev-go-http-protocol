package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/craucrau24/boot-dev-go-http-protocol/internal/response"
)

type Server struct {
	listener *net.TCPListener
	isClosed atomic.Bool
}

func Serve(port int) (*Server, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return nil, fmt.Errorf("couldn't resolve address: %w", err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("couldn't start listening: %w", err)
	}
	server := Server {listener: listener}
	go server.listen()

	return &server, nil
}

func (s *Server) Close() error {
	s.isClosed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for !s.isClosed.Load() {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("error with connection: %v", err)
		} else {
			go s.handle(conn)
		}
	}
}

func (s* Server) handle(conn net.Conn) {
	defer conn.Close()
	response.WriteStatusLine(conn, response.StatusOk)
	response.WriteHeaders(conn, response.GetDefaultHeaders(0))
}