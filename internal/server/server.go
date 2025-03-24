package server

import (
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/VladanT3/TCP_to_HTTP/internal/response"
)

type Server struct {
	Listener net.Listener
	Running  atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		Listener: listener,
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
	err := response.WriteStatusLine(conn, 200)
	if err != nil {
		log.Println(err)
	}

	err = response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	if err != nil {
		log.Println(err)
	}

	conn.Close()
}
