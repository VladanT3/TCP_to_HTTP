package request_line

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLineParse(t *testing.T) {
	// Test: Good GET Request line
	req_line, n, err := ParseRequestLine([]byte("GET / HTTP/1.1\r\n"))
	require.NoError(t, err)
	require.NotNil(t, req_line)
	assert.Equal(t, "GET", req_line.Method)
	assert.Equal(t, "/", req_line.RequestTarget)
	assert.Equal(t, "HTTP/1.1", req_line.HttpVersion)
	assert.Equal(t, 16, n)

	// Test: Good GET Request line with path
	req_line, n, err = ParseRequestLine([]byte("GET /yuh HTTP/1.1\r\n"))
	require.NoError(t, err)
	require.NotNil(t, req_line)
	assert.Equal(t, "GET", req_line.Method)
	assert.Equal(t, "/yuh", req_line.RequestTarget)
	assert.Equal(t, "HTTP/1.1", req_line.HttpVersion)
	assert.Equal(t, 19, n)

	// Test: Good POST Request line with path
	req_line, n, err = ParseRequestLine([]byte("POST /bleeeh HTTP/1.1\r\n"))
	require.NoError(t, err)
	require.NotNil(t, req_line)
	assert.Equal(t, "POST", req_line.Method)
	assert.Equal(t, "/bleeeh", req_line.RequestTarget)
	assert.Equal(t, "HTTP/1.1", req_line.HttpVersion)
	assert.Equal(t, 23, n)

	// Test: Invalid number of parts in request line
	_, _, err = ParseRequestLine([]byte("/bleeeh HTTP/1.1\r\n"))
	require.Error(t, err)

	// Test: Invalid method in request line
	_, _, err = ParseRequestLine([]byte("BLEEEH /yuh HTTP/1.1\r\n"))
	require.Error(t, err)

	// Test: Invalid version in request line
	_, _, err = ParseRequestLine([]byte("PATCH /yuh/bleeeh HTTP/2\r\n"))
	require.Error(t, err)

	// Test: Invalid path in request line
	_, _, err = ParseRequestLine([]byte("PUT /yuh/bl eeeh HTTP/1.1\r\n"))
	require.Error(t, err)
}
