package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:idx], []byte(":"), 2)
	key := string(parts[0])
	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("error: invalid field-name found (%s)", key)
	}
	value := bytes.TrimSpace(parts[1])
	key = strings.TrimSpace(key)
	h.Set(key, string(value))

	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	h[key] = value
}
