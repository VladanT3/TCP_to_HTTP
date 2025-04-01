package request_line

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestLineParsing(t *testing.T) {
	// TEST: Valid Request Line
	data := "GET / HTTP/1.1"
	req_line, err := ParseRequestLine([]byte(data))
	assert.NoError(t, err)
	assert.NotEqual(t, RequestLine{}, req_line)
	assert.Equal(t, "GET", req_line.Method)
	assert.Equal(t, "/", req_line.Target)
	assert.Equal(t, "HTTP/1.1", req_line.HTTPVersion)

	// TEST: Valid POST Request Line
	data = "POST / HTTP/1.1"
	req_line, err = ParseRequestLine([]byte(data))
	assert.NoError(t, err)
	assert.NotEqual(t, RequestLine{}, req_line)
	assert.Equal(t, "POST", req_line.Method)
	assert.Equal(t, "/", req_line.Target)
	assert.Equal(t, "HTTP/1.1", req_line.HTTPVersion)

	// TEST: Valid Request Target Request Line
	data = "PUT /bleeeh/yuuuh HTTP/1.1"
	req_line, err = ParseRequestLine([]byte(data))
	assert.NoError(t, err)
	assert.NotEqual(t, RequestLine{}, req_line)
	assert.Equal(t, "PUT", req_line.Method)
	assert.Equal(t, "/bleeeh/yuuuh", req_line.Target)
	assert.Equal(t, "HTTP/1.1", req_line.HTTPVersion)

	// TEST: Malformed Request Line
	data = "PATCH /blee eh/yuuuh HT TP/1"
	req_line, err = ParseRequestLine([]byte(data))
	assert.Error(t, err)
	assert.Equal(t, RequestLine{}, req_line)

	// TEST: Invalid Method
	data = "YUH /bleeeh/ HTTP/1.1"
	req_line, err = ParseRequestLine([]byte(data))
	assert.Error(t, err)
	assert.Equal(t, RequestLine{}, req_line)

	// TEST: Invalid Request Target
	data = "DELETE bleeeh HTTP/1.1"
	req_line, err = ParseRequestLine([]byte(data))
	assert.Error(t, err)
	assert.Equal(t, RequestLine{}, req_line)

	// TEST: Invalid HTTP Version
	data = "HEAD /bleeeh HTTP/2"
	req_line, err = ParseRequestLine([]byte(data))
	assert.Error(t, err)
	assert.Equal(t, RequestLine{}, req_line)
}
