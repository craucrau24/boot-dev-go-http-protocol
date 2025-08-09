package server

import (
	"fmt"

	"github.com/craucrau24/boot-dev-go-http-protocol/internal/request"
	"github.com/craucrau24/boot-dev-go-http-protocol/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

func HTMLTemplate(title, heading, message string ) string {
	return fmt.Sprintf(`<html>
	<head>
		<title>%s</title>
	</head>
	<body>
		<h1>%s</h1>
		<p>%s</p>
	</body>
	</html>`, title, heading, message)
}

func GetDefaultMessage(status response.StatusCode) (string, string) {
	switch status {
	case response.StatusOk:
		return "200 OK", "Success!"
	case response.StatusBadRequest:
		return "400 Bad Request", "Bad Request"
	case response.StatusInternalServerError:
		return "400 Internal Server Error", "Internal Server Error"
	}

	return "", ""
}