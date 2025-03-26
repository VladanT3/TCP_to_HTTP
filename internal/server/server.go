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

	r, err := request.RequestFromReader(conn)
	if err != nil {
		h_err := &HandlerError{
			StatusCode: response.BadRequest,
			Message:    err.Error(),
		}
		h_err.writeError(conn)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	handler_err := s.Handler(buf, r)
	if handler_err != nil {
		handler_err.writeError(conn)
		return
	}

	b := buf.Bytes()
	headers := response.GetDefaultHeaders(len(b))

	err = response.WriteStatusLine(conn, 200)
	if err != nil {
		log.Println("Error writing response status line:", err)
		return
	}

	err = response.WriteHeaders(conn, headers)
	if err != nil {
		log.Println("Error writing response headers:", err)
		return
	}

	_, err = conn.Write(b)
	if err != nil {
		log.Println("Error writing response body:", err)
	}
}

func (h_err HandlerError) writeError(w io.Writer) {
	err := response.WriteStatusLine(w, h_err.StatusCode)
	if err != nil {
		log.Println("Error writing response status line:", err)
		return
	}

	headers := response.GetDefaultHeaders(len(h_err.Message))
	err = response.WriteHeaders(w, headers)
	if err != nil {
		log.Println("Error writing response headers:", err)
		return
	}

	_, err = w.Write([]byte(h_err.Message))
	if err != nil {
		log.Println("Error writing response body:", err)
		return
	}
}
