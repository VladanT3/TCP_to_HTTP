package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/VladanT3/TCP_to_HTTP/internal/request"
	"github.com/VladanT3/TCP_to_HTTP/internal/server"
)

func handler(w io.Writer, req *request.Request) *server.HandlerError {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		return &server.HandlerError{StatusCode: 400, Message: "Your problem is not my problem\n"}
	} else if req.RequestLine.RequestTarget == "/myproblem" {
		return &server.HandlerError{StatusCode: 500, Message: "Woopsie, my bad\n"}
	} else {
		w.Write([]byte("All good, frfr\n"))
		return nil
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
