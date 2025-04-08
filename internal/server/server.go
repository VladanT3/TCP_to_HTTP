package server

import (
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/VladanT3/TCP_to_HTTP/internal/request"
	"github.com/VladanT3/TCP_to_HTTP/internal/request/headers"
	"github.com/VladanT3/TCP_to_HTTP/internal/response"
)

type Server struct {
	Listener net.Listener
	Handlers map[string]Handler
	Running  atomic.Bool
}

type HTTPError struct {
	StatusCode int
	Message    string
}

type Handler func(w *response.Response, req *request.Request) (string, *HTTPError)

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		Listener: listener,
		Handlers: make(map[string]Handler),
	}
	s.Running.Store(true)

	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	s.Running.Store(false)
	err := s.Listener.Close()
	return err
}

func (s *Server) listen() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			if !s.Running.Load() {
				return
			}
			log.Println("Unable to accept connection.")
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	res := response.Response{
		Data:     conn,
		Headers:  response.GetDefaultHeaders(),
		Trailers: make(headers.Headers),
	}

	req, err := request.ParseRequest(conn)
	if err != nil {
		http_err := HTTPError{StatusCode: 400, Message: ""}
		writeError(&res, &http_err)
		return
	}

	handler, ok := s.Handlers[req.RequestLine.Method+" "+req.RequestLine.Target]
	if !ok {
		http_err := HTTPError{StatusCode: 404, Message: ""}
		writeError(&res, &http_err)
		return
	}

	body, http_err := handler(&res, &req)
	if http_err != nil {
		writeError(&res, http_err)
		return
	}

	res.WriteStatusLine(200)
	res.WriteHeaders(res.Headers)

	_, ok = res.Headers["Chunked-Encoding"]
	if ok {
		for range 10 {
			res.WriteChunkedBody([]byte(body))
		}
		res.WriteChunkedBodyDone()

		_, ok = res.Headers["Trailers"]
		if ok {
			res.WriteTrailers(res.Trailers)
		}
		return
	}

	res.WriteBody([]byte(body))
}

func (s *Server) MapHandler(method string, url string, handler Handler) {
	s.Handlers[method+" "+url] = handler
}

func writeError(res *response.Response, http_err *HTTPError) {
	res.WriteStatusLine(http_err.StatusCode)
	res.Headers["Content-Type"] = "text/plain"
	res.Headers["Content-Length"] = strconv.Itoa(len(http_err.Message))
	res.WriteHeaders(res.Headers)
	res.WriteBody([]byte(http_err.Message))
}
