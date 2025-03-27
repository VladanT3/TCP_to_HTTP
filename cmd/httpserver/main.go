package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/VladanT3/TCP_to_HTTP/internal/request"
	"github.com/VladanT3/TCP_to_HTTP/internal/server"
)

func handler(w io.Writer, req *request.Request) *server.HandlerError {
	fmt.Println("\tENTERED HANDLER FUNC\n\t================================================\n")
	if req.RequestLine.RequestTarget == "/yourproblem" {
		fmt.Println("\tYour problem")
		return &server.HandlerError{StatusCode: 400, Message: "Your problem is not my problem\n"}
	} else if req.RequestLine.RequestTarget == "/myproblem" {
		fmt.Println("\tMy problem")
		return &server.HandlerError{StatusCode: 500, Message: "Woopsie, my bad\n"}
	} else {
		fmt.Println("\tGood request")
		test := bytes.NewBuffer([]byte{})
		multi_writer := io.MultiWriter(w, test)
		_, err := multi_writer.Write([]byte("All good, frfr\n"))
		if err != nil {
			fmt.Println("\tError in handler when writing.")
			return &server.HandlerError{StatusCode: 500, Message: err.Error()}
		}
		fmt.Println("\tWrote successfully in handler.")
		fmt.Printf("\tWrote: %s\n", test.Bytes())
	}
	fmt.Println("\tEXITING HANDLER FUNC\n\t=====================================\n")
	return nil
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
