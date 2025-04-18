package main

import (
	"fmt"
	"log"
	"net"

	"github.com/VladanT3/TCP_to_HTTP/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("Error creating TCP listener: ", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Error accepting connection: ", err)
		}
		fmt.Println("Connection accepted.")

		req, err := request.ParseRequest(conn)
		if err != nil {
			log.Fatal("Unable to proccess request: ", err)
		}

		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", req.RequestLine.Method, req.RequestLine.Target, req.RequestLine.HTTPVersion)
		fmt.Println("Headers:")
		for key, val := range req.Headers {
			fmt.Printf("- %s: %s\n", key, val)
		}
		fmt.Printf("Body:\n%s\n", req.Body)

		fmt.Println("Connection closed.")
	}
}
