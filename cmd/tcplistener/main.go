package main

import (
	"fmt"
	"log"
	"net"

	"github.com/craucrau24/boot-dev-go-http-protocol/internal/request"
)

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
		// fmt.Println("Connection accepted.")
		req, err := request.RequestFromReader(conn)
		fmt.Printf("Request line:\n- Method: %v\n- Target: %v\n- Version: %v\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
		// fmt.Println("Connection closed.")
	}

	//lines := getLinesChannel(file)
	//for line := range lines {
		//fmt.Printf("read: %s\n", line)
	//}
}