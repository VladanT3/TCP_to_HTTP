package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	if string(data[:2]) == "\r\n" {
		return 2, true, nil
	}

	clrf_check := strings.Split(string(data), "\r\n")
	if len(clrf_check) < 2 {
		return 0, false, nil
	}

	split_header := strings.Split(clrf_check[0], ":")
	if split_header[0][len(split_header[0])-1] == ' ' {
		return 0, false, errors.New("Invalid spacing in header.")
	}

	if len(split_header) < 2 {
		return 0, false, errors.New("Malformed header field.")
	}

	key := strings.TrimSpace(split_header[0])
	value := strings.TrimSpace(split_header[1])
	if len(split_header) > 2 {
		split_header[len(split_header)-1] = strings.TrimSpace(split_header[len(split_header)-1])
		for _, val := range split_header[2:] {
			value += ":" + val
		}
	}

	if !isKeyValid(key) {
		return 0, false, errors.New("Invalid character in header key.")
	}

	key = strings.ToLower(key)

	n += len(clrf_check[0]) + 2 // plus 2 for \r and \n because they count as well

	val, ok := h[key]
	if ok {
		h[key] = val + ", " + value
	} else {
		h[key] = value
	}

	return n, false, nil
}

func isKeyValid(key string) bool {
	for _, char := range key {
		if (int(char) >= 48 && int(char) <= 57) || (int(char) >= 65 && int(char) <= 90) || (int(char) >= 97 && int(char) <= 122) {
			continue
		}

		switch true {
		case char == '!':
		case char == '#':
		case char == '$':
		case char == '%':
		case int(char) == 39: // singe tick: '
		case char == '*':
		case char == '+':
		case char == '-':
		case char == '.':
		case char == '^':
		case char == '_':
		case char == '`':
		case char == '|':
		case char == '~':
		default:
			return false
		}
	}

	return true
}

func (h Headers) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	value, ok := h[key]
	if !ok {
		return "", false
	} else {
		return value, true
	}
}
