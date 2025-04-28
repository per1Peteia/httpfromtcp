package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/per1Peteia/httpfromtcp/internal/request"
	"github.com/per1Peteia/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	handler := func(w io.Writer, r *request.Request) *server.HandlerError {
		if r.RequestLine.RequestTarget == "/yourproblem" {
			return server.NewHandlerErr(400, "Your problem is not my problem\n")
		}
		if r.RequestLine.RequestTarget == "/myproblem" {
			return server.NewHandlerErr(500, "Woopsie, my bad\n")
		}
		w.Write([]byte("All good, frfr\n"))
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
