package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lineCh := make(chan string)

	go func() {
		defer f.Close()
		defer close(lineCh)

		buf := make([]byte, 8)
		var lineBuild strings.Builder

		for {
			count, err := f.Read(buf)
			
			if err == io.EOF {
				break
			}

			lineBuild.Write(buf[:count])
			lines := strings.Split(lineBuild.String(), "\n")
			for i := range len(lines) - 1 {
				lineCh <- lines[i]
			}
			lineBuild.Reset()
			lineBuild.WriteString(lines[len(lines) - 1])
		}

		lineCh <- lineBuild.String()
	}()

	return lineCh
}

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("Cannot open file messages.txt: %v\n", err)
	}

	lines := getLinesChannel(file)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}