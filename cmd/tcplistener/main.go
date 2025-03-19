package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	chunk := make([]byte, 8)
	line := ""
	line_ch := make(chan string)

	go func() {
		for {
			n, err := f.Read(chunk)
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal("Error reading from connection: ", err)
			}

			parts := strings.Split(string(chunk[:n]), "\n")
			if len(parts) > 1 {
				line += parts[0]
				line_ch <- line
				line = ""
			}
			line += parts[len(parts)-1]
		}
		line_ch <- line
		close(line_ch)
		f.Close()
	}()

	return line_ch
}

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

		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Printf("%s\n", line)
		}
		fmt.Println("Connection closed.")
	}
}
