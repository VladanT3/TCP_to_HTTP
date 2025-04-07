package handlers

import (
	"github.com/VladanT3/TCP_to_HTTP/internal/request"
	"github.com/VladanT3/TCP_to_HTTP/internal/response"
	"github.com/VladanT3/TCP_to_HTTP/internal/server"
)

func HandleHelloWorld(res *response.ResponseWriter, req *request.Request) (string, *server.HTTPError) {
	body := "<html>\n\t<head>\n\t\t<title>Hello World</title>\n\t</head>\n\t<body>\n\t\t<p>Hello World!</p>\n\t</body>\n</html>\n"
	return body, nil
}

func HandleChunkedEncoding(res *response.ResponseWriter, req *request.Request) (string, *server.HTTPError) {
	body := "chunkedly encoded body\n"
	return body, nil
}
