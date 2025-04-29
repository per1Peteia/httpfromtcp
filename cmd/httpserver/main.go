package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/per1Peteia/httpfromtcp/internal/request"
	"github.com/per1Peteia/httpfromtcp/internal/response"
	"github.com/per1Peteia/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	handler := func(w *response.Writer, r *request.Request) *server.HandlerError {
		if r.RequestLine.RequestTarget == "/yourproblem" {
			return server.NewHandlerErr(400, `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
		}
		if r.RequestLine.RequestTarget == "/myproblem" {
			return server.NewHandlerErr(500, `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
		}
		return nil
	}

	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println(" Server gracefully stopped")
}
