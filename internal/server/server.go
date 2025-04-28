package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/per1Peteia/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{listener: listener}
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	err := response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		log.Printf("error writing status line: %v", err)
	}
	err = response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	if err != nil {
		log.Printf("error writing headers: %v", err)
	}

	return
}
