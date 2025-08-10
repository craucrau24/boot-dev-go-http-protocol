package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/craucrau24/boot-dev-go-http-protocol/internal/request"
	"github.com/craucrau24/boot-dev-go-http-protocol/internal/response"
	"github.com/craucrau24/boot-dev-go-http-protocol/internal/server"
)

const port = 42069

func writeResponse(w *response.Writer, status response.StatusCode, message string) {
		title, heading := server.GetDefaultMessage(status)
		body := []byte(server.HTMLTemplate(title, heading, message))
		headers := response.GetDefaultHeaders(len(body))
		headers.Set("Content-Type", "text/html")
		w.WriteStatusLine(status)
		w.WriteHeaders(headers)
		w.WriteBody(body)
}

func handler(w *response.Writer, req *request.Request) {
	switch  {
	case req.RequestLine.RequestTarget == "/yourproblem":
		writeResponse(w, response.StatusBadRequest, "Your request honestly kinda sucked.")

	case req.RequestLine.RequestTarget == "/myproblem":
		writeResponse(w, response.StatusInternalServerError, "Okay, you know what? This one is on me.")
	case strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin"): {
		path := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")
		resp, err := http.Get(fmt.Sprintf("https://httpbin.org%s", path))
		if err != nil {
			writeResponse(w, response.StatusInternalServerError, "error with httpbin.org")
			return
		}
		defer resp.Body.Close()

		buf := make([]byte, 1024)
		headers := response.GetDefaultHeaders(0)
		headers.Unset("Content-Length")
		headers.Set("Transfer-Encoding", "chunked")
		w.WriteStatusLine(response.StatusOk)
		w.WriteHeaders(headers)
		for {
			n, err := resp.Body.Read(buf)
			if err == nil {
				w.WriteChunkedBody(buf[:n])
			} else {
				w.WriteChunkedBodyDone()
				break
			}
		}
	}
	default:
		writeResponse(w, response.StatusOk, "Your request was an absolute banger.")

	}
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