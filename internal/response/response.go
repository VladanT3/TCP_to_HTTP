package response

import (
	"bytes"
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
	fmt.Println("\t\tENTERED WRITE STATUS LINE FUNC\n\t\t=================================\n")
	test := bytes.NewBuffer([]byte{})
	multi_writer := io.MultiWriter(w, test)
	var err error
	switch status_code {
	case 200:
		_, err = multi_writer.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case 400:
		_, err = multi_writer.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case 500:
		_, err = multi_writer.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	default:
		_, err = multi_writer.Write([]byte("HTTP/1.1 " + strconv.Itoa(int(status_code)) + " \r\n"))
	}
	fmt.Println("\t\tWrote status line.")
	fmt.Printf("\t\tWrote: %s\n", test.Bytes())
	fmt.Println("\t\tEXITING WRITE STATUS LINE FUNC\n\t\t==================================\n")

	return err
}

func GetDefaultHeaders(content_len int) headers.Headers {
	fmt.Println("\t\tENTERED GET DEFAULT HEADERS FUNC\n\t\t==================================\n")
	header := make(headers.Headers)
	header["content-length"] = strconv.Itoa(content_len)
	header["connection"] = "close"
	header["content-type"] = "text/plain"

	fmt.Println("\t\tCreated default headers.")
	fmt.Printf("\t\tCreated headers: %v\n", header)
	fmt.Println("\t\tEXITING GET DEFAULT HEADERS\n\t\t=======================================\n")

	return header
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	fmt.Println("\t\tENTERED WRITE HEADERS FUNC\n\t\t=====================================\n")
	test := bytes.NewBuffer([]byte{})
	multi_writer := io.MultiWriter(w, test)
	var data string
	for key, val := range headers {
		data += fmt.Sprintf("%s: %s\r\n", key, val)
		fmt.Printf("\t\tWrote to data: %s: %s\n", key, val)
	}
	data += "\r\n"
	fmt.Printf("\t\tData: %s\n", data)

	_, err := multi_writer.Write([]byte(data))
	fmt.Printf("\t\tWrote to writer: %s\n", test.Bytes())
	fmt.Println("\t\tEXITING WRITE HEADERS FUNC\n\t\t=======================================\n")
	return err
}
