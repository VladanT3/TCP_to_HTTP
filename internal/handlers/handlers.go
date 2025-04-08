package handlers

import (
	"crypto/sha256"
	"strconv"

	"github.com/VladanT3/TCP_to_HTTP/internal/request"
	"github.com/VladanT3/TCP_to_HTTP/internal/response"
	"github.com/VladanT3/TCP_to_HTTP/internal/server"
)

func HandleHelloWorld(res *response.Response, req *request.Request) (string, *server.HTTPError) {
	body := "<html>\n\t<head>\n\t\t<title>Hello World</title>\n\t</head>\n\t<body>\n\t\t<p>Hello World!</p>\n\t</body>\n</html>\n"
	return body, nil
}

// this is meh but I just want chunked encoding to work and not make a seperate server to send responses
func HandleChunkedEncoding(res *response.Response, req *request.Request) (string, *server.HTTPError) {
	res.Headers["Chunked-Encoding"] = "chunked"
	res.Headers["Content-Type"] = "text/plain"
	res.Headers["Trailers"] = "X-Content-SHA256, X-Content-Length"

	body_base := "chunkedly encoded body\n"
	body := ""
	for range 10 {
		body += body_base
	}

	sha := sha256.Sum256([]byte(body))
	res.Trailers["X-Content-SHA256"] = string(sha[:])
	res.Trailers["X-Content-Length"] = strconv.Itoa(len(body))
	return body_base, nil
}
