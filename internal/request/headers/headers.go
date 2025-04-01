package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

func ParseHeaders(data []byte) (Headers, error) {
	headers := make(Headers)

	parts := strings.SplitSeq(string(data), "\r\n")

	for part := range parts {
		if part == "" {
			break
		}
		part = strings.TrimSpace(part)
		header := strings.Split(part, ":")

		key := header[0]
		if key[len(key)-1] == ' ' {
			return nil, errors.New("Invalid Header Formatting")
		}
		key = strings.ToLower(key)

		val := ""
		if len(header) > 2 {
			val = strings.Join(header[1:], "")
		} else {
			val = header[1]
		}
		val = strings.TrimSpace(val)

		headers[key] = val
	}

	return headers, nil
}

func (h Headers) Get(key string) (string, bool) {
	val, ok := h[key]
	if !ok {
		return "", false
	} else {
		return val, true
	}
}
