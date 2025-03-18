package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/per1Peteia/httpfromtcp/internal/headers"
)

const bufSize = 8

type Request struct {
	RequestLine    RequestLine
	Headers        headers.Headers
	Body           []byte
	bodyLengthRead int
	state          requestState
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateDone
	requestStateParsingHeaders
	requestParsingBody
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufSize, bufSize)
	readToIndex := 0
	req := &Request{
		state:   requestStateInitialized,
		Headers: make(headers.Headers),
		Body:    make([]byte, 0),
	}
	for req.state != requestStateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(io.EOF, err) {
				if req.state != requestStateDone {
					return nil, fmt.Errorf(
						"incomplete request: in state %d, read %d bytes on EOF", req.state, numBytesRead,
					)
				}
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead

		numBytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}
	return req, nil
}

func (req *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for req.state != requestStateDone {
		n, err := req.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			break
		}
		totalBytesParsed += n
	}

	return totalBytesParsed, nil
}

func (req *Request) parseSingle(data []byte) (int, error) {
	switch req.state {
	case requestStateInitialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		req.RequestLine = *requestLine
		req.state = requestStateParsingHeaders // after parsing the requestline is done, switch to parsing headers
		return n, nil
	case requestStateParsingHeaders:
		n, done, err := req.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			req.state = requestParsingBody
		}
		return n, nil
	case requestParsingBody:
		value, exists := req.Headers.Get("content-length")
		if !exists {
			req.state = requestStateDone
			return len(data), nil
		}
		contentLength, err := strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Errorf("error: malformed content-length value: %s", value)
		}

		req.Body = append(req.Body, data...)
		req.bodyLengthRead += len(data)
		if req.bodyLengthRead > contentLength {
			return 0, errors.New("error: length read exceeds content-length")
		}
		if req.bodyLengthRead == contentLength {
			req.state = requestStateDone
		}
		return len(data), nil
	case requestStateDone:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("error: unknown state")
	}
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + 2, nil

}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, errors.New("invalid number of parts in request line")
	}

	method := parts[0]
	if !isUpper(method) {
		return nil, errors.New("method string is not uppercase")
	}

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, errors.New("malformed http version")
	}

	if parts[2] != "HTTP/1.1" {
		return nil, errors.New("wrong HTTP version")
	}
	requestLine := RequestLine{
		HttpVersion:   versionParts[1],
		RequestTarget: parts[1],
		Method:        method,
	}

	return &requestLine, nil
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
