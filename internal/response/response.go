package response

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/VladanT3/TCP_to_HTTP/internal/request/headers"
)

type writer_state int

const (
	status_line writer_state = iota
	field_lines
	body
	done
)

type Writer struct {
	Data        []byte
	writerState writer_state
}

type StatusCode int

const (
	OK                  StatusCode = 200
	BadRequest                     = 400
	InternalServerError            = 500
)

func (w *Writer) WriteStatusLine(status_code StatusCode) error {
	if w.writerState != status_line {
		return errors.New("Writing status line out of order. Make sure you are writing the status line before headers and body.")
	}
	switch status_code {
	case 200:
		w.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case 400:
		w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case 500:
		w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	default:
		w.Write([]byte("HTTP/1.1 " + strconv.Itoa(int(status_code)) + " \r\n"))
	}
	w.writerState = field_lines
	return nil
}

func GetDefaultHeaders(content_len int) headers.Headers {
	header := make(headers.Headers)
	header["content-length"] = strconv.Itoa(content_len)
	header["connection"] = "close"
	header["content-type"] = "text/plain"

	return header
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != field_lines {
		return errors.New("Writing headers out of order. Make sure you write the status line first and write headers before the body.")
	}

	data := ""
	for key, val := range headers {
		data += fmt.Sprintf("%s:%s\r\n", key, val)
	}
	data += "\r\n"

	w.Write([]byte(data))
	w.writerState = body
	return nil
}

func (w *Writer) WriteBody(data []byte) error {
	if w.writerState != body {
		return errors.New("Writing body out of order. Make sure you write the status line and headers first.")
	}
	w.Write(data)
	w.writerState = done
	return nil
}

func (w *Writer) Write(data []byte) {
	w.Data = append(w.Data, data...)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	data := fmt.Sprintf("%X\r\n%s\r\n", len(p), p)
	w.Write([]byte(data))
	return len(data), nil
}
func (w *Writer) WriteChunkedBodyDone() (int, error) {
	w.Write([]byte("0\r\n\r\n"))
	return 5, nil
}
