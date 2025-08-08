package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync/atomic"

	"github.com/craucrau24/boot-dev-go-http-protocol/internal/request"
	"github.com/craucrau24/boot-dev-go-http-protocol/internal/response"
)

type Server struct {
	listener *net.TCPListener
	isClosed atomic.Bool
	handler Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return nil, fmt.Errorf("couldn't resolve address: %w", err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("couldn't start listening: %w", err)
	}
	server := Server {listener: listener, handler: handler}
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

func (s* Server) writeHandlerError(w io.Writer, handleErr HandlerError) error {
	err := response.WriteStatusLine(w, handleErr.Status)
	if err != nil {
		return fmt.Errorf("error while sending response: %w", err)
	}

	msg := []byte(handleErr.Message)
	err = response.WriteHeaders(w, response.GetDefaultHeaders(len(msg)))
	if err != nil {
		return fmt.Errorf("error while sending response: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("error while sending response: %w", err)
	}
	return nil
}

func (s* Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		s.writeHandlerError(conn, HandlerError{Status: response.StatusInternalServerError, Message: fmt.Sprintf("%s", err)})
		return
	}
	var buf bytes.Buffer
	handlerErr := s.handler(&buf, req)

	if handlerErr != nil  {
		s.writeHandlerError(conn, *handlerErr)
	} else {
		response.WriteStatusLine(conn, response.StatusOk)
		response.WriteHeaders(conn, response.GetDefaultHeaders(len(buf.Bytes())))

		conn.Write(buf.Bytes())
	}
}