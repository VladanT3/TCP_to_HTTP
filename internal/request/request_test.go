package request

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestParsing(t *testing.T) {
	// TEST: Valid Request
	buf := bytes.NewBuffer([]byte("GET / HTTP/1.1\r\nContent-Length: 7\r\nConnection: close\r\nHost: localhost:42069\r\nContent-Type: text/plain\r\n\r\nbleeeh\n"))
	req, err := ParseRequest(buf)
	require.NoError(t, err)
	require.NotEqual(t, Request{}, req)
	assert.Equal(t, "GET", req.RequestLine.Method)
	assert.Equal(t, "/", req.RequestLine.Target)
	assert.Equal(t, "HTTP/1.1", req.RequestLine.HTTPVersion)
	assert.Equal(t, "7", req.Headers["content-length"])
	assert.Equal(t, "close", req.Headers["connection"])
	assert.Equal(t, "localhost:42069", req.Headers["host"])
	assert.Equal(t, "text/plain", req.Headers["content-type"])
	assert.Equal(t, "bleeeh\n", string(req.Body))

	// TEST: Valid Request no Content-Length
	buf = bytes.NewBuffer([]byte("GET / HTTP/1.1\r\nConnection: close\r\nHost: localhost:42069\r\nContent-Type: text/plain\r\n\r\nbleeeh\n"))
	req, err = ParseRequest(buf)
	require.NoError(t, err)
	require.NotEqual(t, Request{}, req)
	assert.Equal(t, "GET", req.RequestLine.Method)
	assert.Equal(t, "/", req.RequestLine.Target)
	assert.Equal(t, "HTTP/1.1", req.RequestLine.HTTPVersion)
	assert.Equal(t, "close", req.Headers["connection"])
	assert.Equal(t, "localhost:42069", req.Headers["host"])
	assert.Equal(t, "text/plain", req.Headers["content-type"])
	assert.Nil(t, req.Body)

	// TEST: No CRLF in Request
	buf = bytes.NewBuffer([]byte("GET / HTTP/1.1\nConnection: close\rHost: localhost:42069\n\rContent-Type: text/plain\nbleeeh\r"))
	req, err = ParseRequest(buf)
	require.Error(t, err)
	require.Equal(t, Request{}, req)
}
