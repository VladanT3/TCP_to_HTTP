package request

import (
	"errors"
	"io"
	"strconv"

	"github.com/VladanT3/TCP_to_HTTP/internal/headers"
	"github.com/VladanT3/TCP_to_HTTP/internal/request_line"
)

type state int

const (
	initialized state = iota
	parsing_headers
	parsing_body
	done
)

type Request struct {
	RequestLine request_line.RequestLine
	Headers     headers.Headers
	Body        []byte
	ParserState state
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
		RequestLine: request_line.RequestLine{},
		Headers:     make(headers.Headers),
		Body:        []byte{},
		ParserState: initialized,
	}

	for request.ParserState != done {
		if read_to_index >= buffer_size {
			buf, buffer_size = growSlice(buf)
		}

		bytes_read, err := reader.Read(buf[read_to_index:])
		if err == io.EOF {
			err := request.ValidBody()
			if err != nil {
				return nil, err
			}
			request.ParserState = done
			break
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

func (r *Request) parse(data []byte) (int, error) {
	switch r.ParserState {
	case initialized:
		req_line, bytes_read, err := request_line.ParseRequestLine(data)
		if err != nil {
			return 0, err
		} else if bytes_read == 0 {
			return 0, nil
		} else {
			r.RequestLine = req_line
			r.ParserState = parsing_headers
			return bytes_read, nil
		}
	case parsing_headers:
		n, finished_parsing, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if n == 0 && !finished_parsing {
			return 0, nil
		}

		if finished_parsing {
			r.ParserState = parsing_body
		}

		return n, nil
	case parsing_body:
		_, exists := r.Headers.Get("content-length")
		if !exists && len(data) == 0 {
			r.ParserState = done
			return 0, nil
		} else if !exists && len(data) > 0 {
			return 0, errors.New("Content-Length is absent but body exists.")
		}

		r.Body = append(r.Body, data...)
		return len(data), nil
	case done:
		return 0, errors.New("Trying to read data in a done state.")
	default:
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

func (r *Request) ValidBody() error {
	con_len_str, exists := r.Headers.Get("content-length")
	if !exists {
		con_len_str = "0"
	}
	con_len, err := strconv.Atoi(con_len_str)
	if err != nil {
		return err
	}
	if len(r.Body) != con_len {
		return errors.New("Content-Length doesn't match actual body length.")
	}

	return nil
}
