package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("Cannot open file messages.txt: %v\n", err)
	}

	buf := make([]byte, 8)
	for {
		_, err := file.Read(buf)
		
		if err == io.EOF {
			break
		}

		fmt.Printf("read: %s\n", buf)
	}
}