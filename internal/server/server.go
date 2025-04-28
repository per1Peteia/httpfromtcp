package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/per1Peteia/httpfromtcp/internal/request"
	"github.com/per1Peteia/httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	ErrMsg     string
}

type HandlerFunc func(w io.Writer, req *request.Request) *HandlerError

func NewHandlerErr(code response.StatusCode, msg string) *HandlerError {
	return &HandlerError{
		StatusCode: code,
		ErrMsg:     msg,
	}
}

func (he HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, he.StatusCode)
	msgBytes := []byte(he.ErrMsg)
	response.WriteHeaders(w, response.GetDefaultHeaders(len(msgBytes)))
	w.Write(msgBytes)
}

type Server struct {
	listener net.Listener
	handler  HandlerFunc
	closed   atomic.Bool
}

func Serve(port int, handler HandlerFunc) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{listener: listener, handler: handler}
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

	request, err := request.RequestFromReader(conn)
	if err != nil {
		handlerErr := NewHandlerErr(response.StatusOK, err.Error())
		handlerErr.Write(conn)
		return
	}

	buf := &bytes.Buffer{}

	if handlerErr := s.handler(buf, request); handlerErr != nil {
		handlerErr.Write(conn)
		return
	}

	b := buf.Bytes()
	err = response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		log.Printf("error writing status line: %v", err)
	}
	err = response.WriteHeaders(conn, response.GetDefaultHeaders(len(b)))
	if err != nil {
		log.Printf("error writing headers: %v", err)
	}
	conn.Write(b)
	return

}
