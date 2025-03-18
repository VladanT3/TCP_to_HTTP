package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	chunk := make([]byte, 8)
	line := ""
	line_ch := make(chan string)

	go func() {
		for {
			_, err := f.Read(chunk)
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal("Error reading from file: ", err)
			}

			parts := bytes.Split(chunk, []byte{'\n'})
			line += string(parts[0])
			if len(parts) > 1 {
				line_ch <- line
				line = "" + string(parts[1])
			}
		}
		close(line_ch)
		f.Close()
	}()

	return line_ch
}

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("Unable to open file: ", err)
	}
	lines := getLinesChannel(file)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}
