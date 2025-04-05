package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/VladanT3/TCP_to_HTTP/internal/request"
	"github.com/VladanT3/TCP_to_HTTP/internal/response"
	"github.com/VladanT3/TCP_to_HTTP/internal/server"
)

func handler(res *response.Writer, req *request.Request) {
	if req.RequestLine.Target == "/yourproblem" {
		body := "<html>\n\t<head>\n\t\t<title>400 Bad Request</title>\n\t</head>\n\t<body>\n\t\t<h1>Bad Request</h1>\n\t\t<p>Your request honestly kinda sucked.</p>\n\t</body>\n</html>\n"

		err := res.WriteStatusLine(400)
		if err != nil {
			log.Println("Error writing status line:", err)
			return
		}

		custom_headers := response.GetDefaultHeaders(len(body))
		custom_headers["content-type"] = "text/html"
		err = res.WriteHeaders(custom_headers)
		if err != nil {
			log.Println("Error writing headers:", err)
			return
		}

		err = res.WriteBody([]byte(body))
		if err != nil {
			log.Println("Error writing body:", err)
			return
		}
	} else if req.RequestLine.Target == "/myproblem" {
		body := "<html>\n\t<head>\n\t\t<title>500 Internal Server Error</title>\n\t</head>\n\t<body>\n\t\t<h1>Internal Server Error</h1>\n\t\t<p>Okay, you know what? This one is on me.</p>\n\t</body>\n</html>\n"

		err := res.WriteStatusLine(500)
		if err != nil {
			log.Println("Error writing status line:", err)
			return
		}

		custom_headers := response.GetDefaultHeaders(len(body))
		custom_headers["content-type"] = "text/html"
		err = res.WriteHeaders(custom_headers)
		if err != nil {
			log.Println("Error writing headers:", err)
			return
		}

		err = res.WriteBody([]byte(body))
		if err != nil {
			log.Println("Error writing body:", err)
			return
		}
	} else if strings.HasPrefix(req.RequestLine.Target, "/httpbin/") {
		path := strings.TrimPrefix(req.RequestLine.Target, "/httpbin/")

		err := res.WriteStatusLine(200)
		if err != nil {
			log.Println("Error writing status line:", err)
			return
		}

		custom_headers := response.GetDefaultHeaders(0)
		custom_headers["content-type"] = "application/json"
		custom_headers["transfer-encoding"] = "chunked"
		delete(custom_headers, "content-length")
		err = res.WriteHeaders(custom_headers)
		if err != nil {
			log.Println("Error writing headers:", err)
			return
		}

		resp, err := http.Get("https://httpbin.org/" + path)
		if err != nil {
			log.Println("Error making request to httpbin.org:", err)
			return
		}

		for {
			buf := make([]byte, 1024)
			n, err := resp.Body.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal("Error reading body from httpbin.org:", err)
			}

			_, err = res.WriteChunkedBody(buf[:n])
			if err != nil {
				log.Println("Error writing chunked body:", err)
				return
			}
		}

		_, err = res.WriteChunkedBodyDone()
		if err != nil {
			log.Println("Error writing chunked body done:", err)
			return
		}
	} else {
		body := "<html>\n\t<head>\n\t\t<title>200 OK</title>\n\t</head>\n\t<body>\n\t\t<h1>Success!</h1>\n\t\t<p>Your request was an absolute banger.</p>\n\t</body>\n</html>\n"

		err := res.WriteStatusLine(200)
		if err != nil {
			log.Println("Error writing status line:", err)
			return
		}

		custom_headers := response.GetDefaultHeaders(len(body))
		custom_headers["content-type"] = "text/html"
		err = res.WriteHeaders(custom_headers)
		if err != nil {
			log.Println("Error writing headers:", err)
			return
		}

		err = res.WriteBody([]byte(body))
		if err != nil {
			log.Println("Error writing body:", err)
			return
		}
	}
}

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}

	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	server.Close()
	log.Println("Server stopped.")
}
