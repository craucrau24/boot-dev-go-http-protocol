package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("Cannot open file messages.txt: %v\n", err)
	}

	buf := make([]byte, 8)
	var lineBuild strings.Builder

	for {
		count, err := file.Read(buf)
		
		if err == io.EOF {
			break
		}

		lineBuild.Write(buf[:count])
		lines := strings.Split(lineBuild.String(), "\n")
		for i := range len(lines) - 1 {
			fmt.Printf("read: %s\n", lines[i])
		}
		lineBuild.Reset()
		lineBuild.WriteString(lines[len(lines) - 1])
	}
	fmt.Printf("read: %s\n", &lineBuild)
}