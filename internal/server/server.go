package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/VladanT3/TCP_to_HTTP/internal/request"
	"github.com/VladanT3/TCP_to_HTTP/internal/response"
)

type Server struct {
	Listener net.Listener
	Handler  Handler
	Running  atomic.Bool
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		Listener: listener,
		Handler:  handler,
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

	req, err := request.RequestFromReader(conn)
	if err != nil {
		h_err := &HandlerError{
			StatusCode: 400,
			Message:    "Bad Request",
		}
		h_err.writeError(conn)
		return
	}

	res_body := bytes.NewBuffer([]byte{})
	h_err := s.Handler(res_body, req)
	if h_err != nil {
		h_err.writeError(conn)
		return
	}

	res := bytes.NewBuffer([]byte{})

	err = response.WriteStatusLine(res, 200)
	if err != nil {
		log.Println("Error writing response status line:", err)
		return
	}

	default_headers := response.GetDefaultHeaders(len(res_body.Bytes()))
	err = response.WriteHeaders(res, default_headers)
	if err != nil {
		log.Println("Error writing response headers:", err)
		return
	}

	_, err = res.Write(res_body.Bytes())
	if err != nil {
		log.Println("Error writing response body:", err)
		return
	}

	_, err = conn.Write(res.Bytes())
	if err != nil {
		log.Println("Error writing response:", err)
		return
	}
}

func (h_err HandlerError) writeError(w io.Writer) {
	res := bytes.NewBuffer([]byte{})

	err := response.WriteStatusLine(res, h_err.StatusCode)
	if err != nil {
		log.Println("Error writing error response status line:", err)
		return
	}

	default_headers := response.GetDefaultHeaders(len(h_err.Message))
	err = response.WriteHeaders(res, default_headers)
	if err != nil {
		log.Println("Error writing error response headers:", err)
		return
	}

	_, err = res.Write([]byte(h_err.Message))
	if err != nil {
		log.Println("Error writing error response body:", err)
		return
	}

	_, err = w.Write(res.Bytes())
	if err != nil {
		log.Println("Error writing error response:", err)
		return
	}
}
