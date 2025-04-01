package body

import (
	"errors"
	"strconv"
)

func ParseBody(data []byte, con_len_str string) ([]byte, error) {
	con_len, err := strconv.Atoi(con_len_str)
	if err != nil {
		return nil, errors.New("Invalid Content-Length Header")
	}

	if len(data) != con_len {
		return nil, errors.New("Invalid Body Length")
	}

	return data, nil
}
