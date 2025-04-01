package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParsing(t *testing.T) {
	// TEST: Valid Headers
	data := "Content-Length: 15\r\nConnection: close\r\nContent-Type: text/plain\r\nHost: localhost:42069\r\n\r\n"
	headers, err := ParseHeaders([]byte(data))
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "15", headers["content-length"])
	assert.Equal(t, "close", headers["connection"])
	assert.Equal(t, "text/plain", headers["content-type"])
	assert.Equal(t, "localhost:42069", headers["host"])

	// TEST: Valid Headers With a lot of spaces
	data = "        Content-Length: 15   \r\n    Connection:   close  \r\n    Content-Type:text/plain    \r\n   Host:    localhost:42069   \r\n\r\n"
	headers, err = ParseHeaders([]byte(data))
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "15", headers["content-length"])
	assert.Equal(t, "close", headers["connection"])
	assert.Equal(t, "text/plain", headers["content-type"])
	assert.Equal(t, "localhost:42069", headers["host"])

	// TEST: Valid small Header
	data = "Content-Length: 15\r\n\r\n"
	headers, err = ParseHeaders([]byte(data))
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "15", headers["content-length"])

	// TEST: No Header
	data = "\r\n"
	headers, err = ParseHeaders([]byte(data))
	require.NoError(t, err)
	require.NotNil(t, headers)

	// TEST: Invalid small Header
	data = "Content-Length: 15\r\n"
	headers, err = ParseHeaders([]byte(data))
	require.Error(t, err)
	require.Nil(t, headers)

	// TEST: Malformed Headers
	data = "Content-Length: 15\r\nConnection: close\r\nContent-Type: text/plain\r\n"
	headers, err = ParseHeaders([]byte(data))
	assert.Error(t, err)
	assert.Nil(t, headers)

	// TEST: Invalid Header Formatting
	data = "Content-Length  : 15\r\nConnection : close\r\nContent-Type    : text/plain\r\n\r\n"
	headers, err = ParseHeaders([]byte(data))
	assert.Error(t, err)
	assert.Nil(t, headers)
}
