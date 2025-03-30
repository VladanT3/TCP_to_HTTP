package server

import (
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

type Handler func(w *response.Writer, req *request.Request)

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

	res := response.Writer{}
	req, err := request.RequestFromReader(conn)
	if err != nil {
		err := res.WriteStatusLine(400)
		if err != nil {
			log.Println("Error writing status line:", err)
			return
		}
		err = res.WriteHeaders(nil)
		if err != nil {
			log.Println("Error writing headers:", err)
			return
		}

		_, err = conn.Write(res.Data)
		if err != nil {
			log.Println("Error writing to connection:", err)
			return
		}
		return
	}

	s.Handler(&res, req)
	_, err = conn.Write(res.Data)
	if err != nil {
		log.Println("Error writing to connection:", err)
		return
	}
}
