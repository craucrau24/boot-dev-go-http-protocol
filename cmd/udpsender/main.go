package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udp, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalln("Couldn't resolve localhost:42069")
	}
	conn, err := net.DialUDP("udp", nil, udp)
	if err != nil {
		log.Fatalf("Couldn't initiate UDP network: %v", err)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln("Error reading from standard input")
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Fatalln("Error writing to UDP network")
		}
	}
}