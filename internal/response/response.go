package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/VladanT3/TCP_to_HTTP/internal/headers"
)

type StatusCode int

const (
	OK                  StatusCode = 200
	BadRequest                     = 400
	InternalServerError            = 500
)

func WriteStatusLine(w io.Writer, status_code StatusCode) error {
	var err error
	switch status_code {
	case 200:
		_, err = w.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case 400:
		_, err = w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case 500:
		_, err = w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	default:
		_, err = w.Write([]byte("HTTP/1.1 " + strconv.Itoa(int(status_code)) + " \r\n"))
	}

	return err
}

func GetDefaultHeaders(content_len int) headers.Headers {
	header := make(headers.Headers)
	header["Content-Length"] = strconv.Itoa(content_len)
	header["Connection"] = "close"
	header["Content-Type"] = "text/plain"

	return header
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	var data string
	for key, val := range headers {
		data += fmt.Sprintf("%s: %s\r\n", key, val)
	}
	data += "\r\n"

	_, err := w.Write([]byte(data))
	return err
}
