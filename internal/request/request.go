package request

import (
	"errors"
	"io"
	"strings"

	"github.com/VladanT3/TCP_to_HTTP/internal/request/body"
	"github.com/VladanT3/TCP_to_HTTP/internal/request/headers"
	"github.com/VladanT3/TCP_to_HTTP/internal/request/request_line"
)

type Request struct {
	RequestLine request_line.RequestLine
	Headers     headers.Headers
	Body        []byte
}

func ParseRequest(conn io.Reader) (Request, error) {
	req_buf := make([]byte, 1024)
	total := 0
	var err error

	for {
		n, err := conn.Read(req_buf[total:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return Request{}, err
		}
		if n < 1 {
			break
		}
		total += n
	}

	request := Request{}

	pos := strings.Index(string(req_buf), "\r\n")
	if pos == -1 {
		return Request{}, errors.New("Malformed Request")
	}
	body_pos := strings.LastIndex(string(req_buf), "\r\n")

	request.RequestLine, err = request_line.ParseRequestLine(req_buf[:pos])
	if err != nil {
		return Request{}, err
	}

	request.Headers, err = headers.ParseHeaders(req_buf[pos+2 : body_pos+2])
	if err != nil {
		return Request{}, err
	}

	con_len, exists := request.Headers.Get("content-length")
	if !exists {
		return request, nil
	}
	request.Body, err = body.ParseBody(req_buf[body_pos+2:total], con_len)
	if err != nil {
		return Request{}, err
	}

	return request, nil
}
