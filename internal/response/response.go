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
	_, err := w.Write(getStatusLine(statusCode))
	return err
}

func getStatusLine(statusCode StatusCode) []byte {
	var reasonPhrase string
	switch statusCode {
	case StatusOK:
		reasonPhrase = "OK"
	case StatusBadRequest:
		reasonPhrase = "Bad Request"
	case StatusInternalServerError:
		reasonPhrase = "Internal Server Error"
	}
	return []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase))
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
