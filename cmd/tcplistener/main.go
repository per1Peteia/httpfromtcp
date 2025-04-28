package main

import (
	"fmt"
	"github.com/per1Peteia/httpfromtcp/internal/request"
	"log"
	"net"
)

const port string = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening: %s", err.Error())
	}
	defer listener.Close()
	fmt.Println("listening on port:", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error accepting: %s", err.Error())
			continue
		}
		fmt.Printf("connection accepted from: %s\n", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error requesting: %s", err.Error())
		}
		fmt.Printf(
			"Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n",
			req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion,
		)
		fmt.Println("Headers:")
		for key, value := range req.Headers {
			fmt.Printf("- %v: %v\n", key, value)
		}
		fmt.Printf("Body:\n%s\n", req.Body)
	}
}
