package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

// takes a bytestring (request headers) and returns bytes consumed, parsing status and error status
func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// looks if it has enough sent data
	idx := bytes.Index(data, []byte("\r\n"))
	// if not, need more
	if idx == -1 {
		return 0, false, nil
	}
	// look for end of headers, if found: parsing done
	if idx == 0 {
		return 2, true, nil
	}

	// parse the data
	parts := bytes.SplitN(data[:idx], []byte(":"), 2)
	key := strings.ToLower(string(parts[0]))
	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("error: invalid field-name found (%s) - whitespace", key)
	}
	key = strings.TrimSpace(key)
	if !validFieldName(key) {
		return 0, false, fmt.Errorf("error: invalid field-name (%s) - illegal codepoint", key)
	}
	value := bytes.TrimSpace(parts[1])

	h.Set(key, string(value))

	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	if _, exists := h[key]; exists {
		h[key] += ", " + value
		return
	}
	h[key] = value
}

// range table to use as valid codepoints for field-names
var validRunes = &unicode.RangeTable{
	R16: []unicode.Range16{
		{
			Lo:     0x0021,
			Hi:     0x007E,
			Stride: 1,
		},
	},
}

// takes a string as input and checks whether it is within the range of valid, specified unicode codepoints
func validFieldName(str string) bool {
	for _, char := range str {
		if !unicode.In(char, validRunes) {
			return false
		}
	}
	return true
}
