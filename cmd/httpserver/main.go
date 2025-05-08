package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/per1Peteia/httpfromtcp/internal/headers"
	"github.com/per1Peteia/httpfromtcp/internal/request"
	"github.com/per1Peteia/httpfromtcp/internal/response"
	"github.com/per1Peteia/httpfromtcp/internal/server"
)

const port = 42069

func main() {
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

func handler(w *response.Writer, r *request.Request) {
	if r.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, r)
		return
	}
	if r.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, r)
		return
	}
	if strings.HasPrefix(r.RequestLine.RequestTarget, "/httpbin/") {
		proxyHandler(w, r)
		return
	}
	if r.RequestLine.RequestTarget == "/video" {
		videoHandler(w, r)
		return
	}

	handler200(w, r)
	return
}

func videoHandler(w *response.Writer, r *request.Request) {
	w.WriteStatusLine(response.StatusOK)
	file, err := os.ReadFile("assets/vim.mp4")
	if err != nil {
		handler500(w, r)
		return
	}
	h := response.GetDefaultHeaders(len(file))
	h.Reset("Content-Type", "video/mp4")

	w.WriteHeaders(h)
	w.WriteBody(file)
	return
}

func proxyHandler(w *response.Writer, r *request.Request) {
	endpoint := strings.TrimPrefix(r.RequestLine.RequestTarget, "/httpbin/")
	res, err := http.Get(fmt.Sprintf("https://httpbin.org/%s", endpoint))
	if err != nil {
		handler500(w, r)
		return
	}
	defer res.Body.Close()

	w.WriteStatusLine(response.StatusOK)
	h := response.GetDefaultHeaders(0)
	h.Reset("Transfer-Encoding", "chunked")
	h.Remove("Content-Length")
	h.Set("Trailer", "X-Content-SHA256")
	h.Set("Trailer", "X-Content-Length")
	w.WriteHeaders(h)

	const maxChunkSize = 1024
	var body []byte
	buf := make([]byte, maxChunkSize)
	for {
		bytesRead, err := res.Body.Read(buf)
		fmt.Println(bytesRead, "bytes read")
		if bytesRead > 0 {
			_, err := w.WriteChunkedBody(buf[:bytesRead])
			if err != nil {
				fmt.Println("error writing chunked body", err)
				break
			}
			body = append(body, buf[:bytesRead]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("error reading response body:", err)
			break
		}
	}
	err = w.WriteChunkedBodyDone(body)
	if err != nil {
		fmt.Println("error writing chunked body done:", err)
	}

	trailers := make(headers.Headers, 2)
	trailers.Set("X-Content-SHA256", fmt.Sprintf("%x", sha256.Sum256(body)))
	trailers.Set("X-Content-Length", fmt.Sprintf("%d", len(body)))
	err = w.WriteTrailers(trailers)
	if err != nil {
		fmt.Println("error writing trailers", err)
	}
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusBadRequest)
	body := []byte(`<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p	>Your request honestly kinda sucked.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Reset("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusInternalServerError)
	body := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Reset("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusOK)
	body := []byte(`<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Reset("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}
