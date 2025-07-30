package main

import (
	"fmt"
	"io"
	"log"
	"net"
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
	const PORT = 42069

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", PORT))
	if err != nil {
		log.Fatalf("Cannot open port %d for listening: %v\n", PORT, err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Cannot accept connection: %v\n", err)
		}
		fmt.Println("Connection accepted.")
		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Println(line)
		}
		fmt.Println("Connection closed.")
	}

	//lines := getLinesChannel(file)
	//for line := range lines {
		//fmt.Printf("read: %s\n", line)
	//}
}