package request

import (
	"errors"
	"io"
	"strings"
)

type state int

const (
	initialized state = iota
	done
)

type Request struct {
	RequestLine RequestLine
	ParserState state
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer_size := 8
	buf := make([]byte, buffer_size, buffer_size)

	read_to_index := 0

	request := &Request{
		ParserState: 0,
	}

	for request.ParserState != 1 {
		if read_to_index >= buffer_size {
			buf, buffer_size = growSlice(buf)
		}

		bytes_read, err := reader.Read(buf[read_to_index:])
		if err == io.EOF {
			request.ParserState = 1
			if bytes_read == 0 {
				break
			}
		} else if err != nil {
			return nil, err
		}
		read_to_index += bytes_read

		bytes_parsed, err := request.parse(buf[:read_to_index])
		if err != nil {
			return nil, err
		}
		if bytes_parsed > 0 {
			copy(buf, buf[bytes_parsed:read_to_index])
			read_to_index -= bytes_parsed
		}
	}

	return request, nil
}

func parseRequestLine(data []byte) (RequestLine, int, error) {
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

	return req_line, len(split_req[0]), nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.ParserState == 0 {
		req_line, bytes_read, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		} else if bytes_read == 0 {
			return 0, nil
		} else {
			r.RequestLine = req_line
			r.ParserState = 1
			return bytes_read, nil
		}
	} else if r.ParserState == 1 {
		return 0, errors.New("Trying to read data in a done state.")
	} else {
		return 0, errors.New("Unknown state.")
	}
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}

	endIndex := min(cr.pos+cr.numBytesPerRead, len(cr.data))

	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n
	if n > cr.numBytesPerRead {
		n = cr.numBytesPerRead
		cr.pos -= n - cr.numBytesPerRead
	}
	return n, nil
}

func growSlice[T any](slice []T) ([]T, int) {
	new_slice := make([]T, cap(slice)*2, cap(slice)*2)
	copy(new_slice, slice)
	return new_slice, cap(new_slice)
}
