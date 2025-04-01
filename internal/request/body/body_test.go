package body

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsingBody(t *testing.T) {
	// TEST: Valid Body
	data := "bleeeh"
	body, err := ParseBody([]byte(data), "6")
	require.NoError(t, err)
	require.NotNil(t, body)
	assert.Equal(t, "bleeeh", string(body))

	// TEST: No Content-Length and Empty Body
	data = ""
	body, err = ParseBody([]byte(data), "0")
	require.NoError(t, err)
	require.NotNil(t, body)
	assert.Equal(t, "", string(body))

	// TEST: Invalid Content-Length Header
	data = "bleeeh"
	body, err = ParseBody([]byte(data), "a")
	require.Error(t, err)
	require.Nil(t, body)

	// TEST: Invalid Body Length
	data = "bleeeh"
	body, err = ParseBody([]byte(data), "4")
	require.Error(t, err)
	require.Nil(t, body)
}
