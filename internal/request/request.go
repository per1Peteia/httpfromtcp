package request

import (
	"errors"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return &Request{}, err
	}

	requestLine := strings.Split(string(bytes), "\r\n")[0]
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return &Request{}, errors.New("invalid number of parts in request line")
	}
	if !isUpper(parts[0]) {
		return &Request{}, errors.New("method string is not uppercase")
	}
	if parts[2] != "HTTP/1.1" {
		return &Request{}, errors.New("wrong HTTP version")
	}
	request := Request{
		RequestLine: RequestLine{
			HttpVersion:   strings.Split(parts[2], "/")[1],
			RequestTarget: parts[1],
			Method:        parts[0],
		},
	}
	return &request, nil
}

func isUpper(str string) bool {
	if len(str) == 0 {
		return false
	}

	for _, char := range str {
		if !unicode.IsUpper(char) {
			return false
		}
	}
	return true
}
