package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/JoStMc/httpfromtcp/internal/request"
	"github.com/JoStMc/httpfromtcp/internal/response"
	"github.com/JoStMc/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	s, err := server.Serve(port, handlerPaths)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer s.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handlerPaths(w *response.Writer, req *request.Request) {
	var b []byte
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		b = []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>your request was bad</p>
  </body>
</html>`)
		w.WriteStatusLine(response.StatusBadRequest)
	case "/myproblem":
		b = []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>This one is on me.</p>
  </body>
</html>`)
		w.WriteStatusLine(response.StatusIntervalServerError)
	default:
		b = []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
		w.WriteStatusLine(response.StatusOK)
	}

	headers := response.GetDefaultHeaders(len(b))
	headers.Replace("Content-Type", "text/html")

	w.WriteHeaders(headers)
	w.WriteBody(b)
} 
