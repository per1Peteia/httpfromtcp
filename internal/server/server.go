package server

import (
	"fmt"
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

type HandlerFunc func(w *response.Writer, req *request.Request) *HandlerError

func NewHandlerErr(code response.StatusCode, msg string) *HandlerError {
	return &HandlerError{
		StatusCode: code,
		ErrMsg:     msg,
	}
}

func (he HandlerError) Write(w *response.Writer) {
	w.WriteStatusLine(he.StatusCode)
	msgBytes := []byte(he.ErrMsg)
	headers := response.GetDefaultHeaders(len(msgBytes))
	headers.Reset("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(msgBytes)
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

	w := &response.Writer{}
	r, err := request.RequestFromReader(conn)
	if err != nil {
		handlerErr := NewHandlerErr(response.StatusOK, err.Error())
		handlerErr.Write(w)
		conn.Write(w.Response)
		return
	}

	// error handling response
	if handlerErr := s.handler(w, r); handlerErr != nil {
		handlerErr.Write(w)
		conn.Write(w.Response)
		return
	}

	// no errors means default response 200 is written to the response.Writer
	body := []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
	err = w.WriteStatusLine(response.StatusOK)
	if err != nil {
		log.Printf("error writing status line: %v", err)
	}
	headers := response.GetDefaultHeaders(len(body))
	headers.Reset("Content-Type", "text/html")
	err = w.WriteHeaders(headers)
	if err != nil {
		log.Printf("error writing headers: %v", err)
	}
	w.WriteBody(body)

	conn.Write(w.Response)
	return
}
