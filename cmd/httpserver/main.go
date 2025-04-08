package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/VladanT3/TCP_to_HTTP/internal/handlers"
	"github.com/VladanT3/TCP_to_HTTP/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}

	server.MapHandler("GET", "/hello/world", handlers.HandleHelloWorld)
	server.MapHandler("GET", "/chunked", handlers.HandleChunkedEncoding)

	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	server.Close()
	log.Println("Server stopped.")
}
