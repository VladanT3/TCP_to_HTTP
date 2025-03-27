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
	log.Println("ENTERED SERVER.HANDLE FUNC\n============================================\n")
	defer conn.Close()

	r, err := request.RequestFromReader(conn)
	log.Printf("Request: %v\n", r)
	if err != nil {
		log.Println("Error parsing request.")
		h_err := &HandlerError{
			StatusCode: response.BadRequest,
			Message:    err.Error(),
		}
		log.Printf("Handler error: %v\n", h_err)
		h_err.writeError(conn)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	handler_err := s.Handler(buf, r)
	if handler_err != nil {
		log.Println("Handler error isnt nil.")
		log.Printf("Handler error: %v\n", handler_err)
		handler_err.writeError(conn)
		return
	}

	b := buf.Bytes()
	log.Printf("What we got back from the handler: %s\n", b)
	log.Printf("Len of b: %d\n", len(b))
	headers := response.GetDefaultHeaders(len(b))
	log.Printf("Headers: %v\n", headers)

	err = response.WriteStatusLine(conn, 200)
	if err != nil {
		log.Println("Error writing response status line:", err)
		return
	}
	log.Println("Wrote status line.")

	err = response.WriteHeaders(conn, headers)
	if err != nil {
		log.Println("Error writing response headers:", err)
		return
	}
	log.Println("Wrote headers.")

	test := bytes.NewBuffer([]byte{})
	multi_writer := io.MultiWriter(conn, test)
	_, err = multi_writer.Write(b)
	if err != nil {
		log.Println("Error writing response body:", err)
	}
	log.Println("Wrote body.")
	log.Printf("Wrote: %s\n", test.Bytes())
	log.Println("EXITING SERVER.HANDLE FUNC\n=======================================\n")
}

func (h_err HandlerError) writeError(w io.Writer) {
	log.Println("\tENTERED WRITE ERROR FUNC\n\t========================================\n")
	err := response.WriteStatusLine(w, h_err.StatusCode)
	if err != nil {
		log.Println("Error writing response status line:", err)
		return
	}
	log.Println("\tWrote status line.")

	headers := response.GetDefaultHeaders(len(h_err.Message))
	err = response.WriteHeaders(w, headers)
	if err != nil {
		log.Println("Error writing response headers:", err)
		return
	}
	log.Println("\tWrote headers.")

	test := bytes.NewBuffer([]byte{})
	multi_writer := io.MultiWriter(w, test)
	_, err = multi_writer.Write([]byte(h_err.Message))
	if err != nil {
		log.Println("Error writing response body:", err)
		return
	}
	log.Println("\tWrote body.")
	log.Printf("\tWrote: %s\n", test.Bytes())
	log.Println("\tEXITING WRITE ERROR FUNC\n\t============================================\n")
}
