package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/per1Peteia/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusContinue StatusCode = (iota + 1) * 100
	StatusOK
	StatusMultipleChoices
	StatusBadRequest
	StatusInternalServerError
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusOK:
		_, err := fmt.Fprint(w, "HTTP/1.1 200 OK\r\n")
		if err != nil {
			return err
		}
		return nil
	case StatusBadRequest:
		_, err := fmt.Fprint(w, "HTTP/1.1 400 Bad Request\r\n")
		if err != nil {
			return err
		}
		return nil
	case StatusInternalServerError:
		_, err := fmt.Fprint(w, "HTTP/1.1 500 Internal Server Error\r\n")
		if err != nil {
			return err
		}
		return nil
	default:
		_, err := fmt.Fprint(w, "HTTP/1.1 \r\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func GetDefaultHeaders(contentLength int) headers.Headers {
	headers := make(headers.Headers, 3)
	headers.Set("Content-Length", strconv.Itoa(contentLength))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := fmt.Fprintf(w, "%s: %s\r\n", key, value)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprint(w, "\r\n")
	if err != nil {
		return err
	}
	return nil
}
