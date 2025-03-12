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

	// refactor this to use strings.Split like a sane person would
	strippedHeaderStr := strings.TrimSpace(string(data[:idx]))
	colIdx := strings.Index(strippedHeaderStr, ":")
	if colIdx != -1 && strippedHeaderStr[colIdx-1] == ' ' {
		return 0, false, fmt.Errorf("error: whitespace suffix to field-name found (%s)", strippedHeaderStr)
	}
	h[strippedHeaderStr[:colIdx]] = strings.TrimSpace(strippedHeaderStr[colIdx+1:])
	return len(data[:idx+2]), false, nil
}
