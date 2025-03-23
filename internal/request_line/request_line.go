package request_line

import (
	"errors"
	"strings"
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func ParseRequestLine(data []byte) (RequestLine, int, error) {
	split_req := strings.Split(string(data), "\r\n")

	if len(split_req) < 2 {
		return RequestLine{}, 0, nil
	}

	req_line_parts := strings.Split(split_req[0], " ")
	if len(req_line_parts) != 3 {
		return RequestLine{}, 0, errors.New("Invalid request line.")
	}

	method := req_line_parts[0]
	target := req_line_parts[1]
	http_ver := req_line_parts[2]
	if http_ver != "HTTP/1.1" {
		return RequestLine{}, 0, errors.New("Invalid HTTP version.")
	}
	if !strings.Contains(target, "/") || strings.Contains(target, " ") {
		return RequestLine{}, 0, errors.New("Invalid request target.")
	}
	switch method {
	case "CONNECT":
	case "DELETE":
	case "GET":
	case "HEAD":
	case "OPTIONS":
	case "PATCH":
	case "POST":
	case "PUT":
	case "TRACE":
	default:
		return RequestLine{}, 0, errors.New("Invalid method.")
	}
	req_line := RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   http_ver,
	}

	return req_line, len(split_req[0]) + 2, nil
}
