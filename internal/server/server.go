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

type Handler func(w *response.ResponseWriter, req *request.Request) (string, *HTTPError)

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
	res := response.ResponseWriter{}
	res_header := make(headers.Headers)

	req, err := request.ParseRequest(conn)
	if err != nil {
		res.WriteResponse(400, res_header, "")
		_, err := conn.Write(res.Data)
		if err != nil {
			log.Println("Error writing response:", err)
		}
		return
	}

	handler, ok := s.Handlers[req.RequestLine.Method+" "+req.RequestLine.Target]
	if !ok {
		res.WriteResponse(404, res_header, "")
		_, err := conn.Write(res.Data)
		if err != nil {
			log.Println("Error writing response:", err)
		}
		return
	}
	body, http_err := handler(&response.ResponseWriter{}, &req)
	if http_err != nil {
		res.WriteResponse(http_err.StatusCode, res_header, http_err.Message)
		_, err := conn.Write(res.Data)
		if err != nil {
			log.Println("Error writing response:", err)
		}
		return
	}

	res.WriteResponse(200, res_header, body)
	_, err = conn.Write(res.Data)
	if err != nil {
		log.Println("Error writing response:", err)
	}
	return
}

func (s *Server) MapHandler(method string, url string, handler Handler) {
	s.Handlers[method+" "+url] = handler
}
