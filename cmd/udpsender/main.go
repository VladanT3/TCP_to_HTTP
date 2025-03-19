package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal("Unable to resolve UDP address: ", err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("Unable to establish connection: ", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Print("Invalid input: ", err, "\n")
		}
		_, err = conn.Write([]byte(input))
		if err != nil {
			log.Print("Error writing to UDP connection: ", err, "\n")
		}
	}
}
