package request_line

import (
	"bytes"
	"errors"
)

type RequestLine struct {
	Method      string
	Target      string
	HTTPVersion string
}

func ParseRequestLine(data []byte) (RequestLine, error) {
	parts := bytes.Split(data, []byte{' '})

	if len(parts) != 3 {
		return RequestLine{}, errors.New("Malformed Request Line")
	}

	method := parts[0]
	target := parts[1]
	httpver := parts[2]

	switch string(method) {
	case "GET":
	case "HEAD":
	case "OPTIONS":
	case "TRACE":
	case "PUT":
	case "DELETE":
	case "POST":
	case "PATCH":
	case "CONNECT":
	default:
		return RequestLine{}, errors.New("Invalid Method")
	}

	if !bytes.Contains(target, []byte{'/'}) {
		return RequestLine{}, errors.New("Invalid Request Target")
	}

	if string(httpver) != "HTTP/1.1" {
		return RequestLine{}, errors.New("Invalid HTTP Version")
	}

	return RequestLine{string(method), string(target), string(httpver)}, nil
}
