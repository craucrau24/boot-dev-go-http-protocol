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
		count, err := file.Read(buf)
		// fmt.Printf("%v, %v", count, err)
		
		if err == io.EOF {
			break
		}

		fmt.Printf("read: %s\n", buf[0:count])
	}
}