package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	unfiltered_req, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	str_unf_req := string(unfiltered_req)
	split_req := strings.Split(str_unf_req, "\r\n")
	req_line_str := split_req[0]
	req_line_parts := strings.Split(req_line_str, " ")
	if len(req_line_parts) != 3 {
		return nil, errors.New("Invalid request line.")
	}
	method := req_line_parts[0]
	target := req_line_parts[1]
	http_ver := req_line_parts[2]
	if http_ver != "HTTP/1.1" {
		return nil, errors.New("Invalid HTTP version.")
	}
	if !strings.Contains(target, "/") || strings.Contains(target, " ") {
		return nil, errors.New("Invalid request target.")
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
		return nil, errors.New("Invalid method.")
	}
	req_line := RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   http_ver,
	}

	return &Request{RequestLine: req_line}, nil
}
