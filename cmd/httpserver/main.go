package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/craucrau24/boot-dev-go-http-protocol/internal/request"
	"github.com/craucrau24/boot-dev-go-http-protocol/internal/response"
	"github.com/craucrau24/boot-dev-go-http-protocol/internal/server"
)

const port = 42069

func handler(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{Status: response.StatusBadRequest, Message: "Your problem is not my problem\n"}
	case "/myproblem":
		return &server.HandlerError{Status: response.StatusInternalServerError, Message: "Woopsie, my bad\n"}
	}
	w.Write([]byte("All good, frfr\n"))
	return nil
}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}