package main

import (
	"fmt"
	"strings"
)

func main() {
	data := "Host: localhost:42069\r\nYuh: bleeeh\r\n\r\n"
	split := strings.Split(data, "\r\n")
	fmt.Println(len(split))
	for _, part := range split {
		fmt.Println(len(part))
		fmt.Println(part)
	}
}
