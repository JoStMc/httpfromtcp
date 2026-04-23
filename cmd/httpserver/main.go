package main

import (
	"io"
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

func handlerPaths(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return server.NewHandlerError(response.StatusBadRequest, []byte("Your problem\n"))
	case "/myproblem":
		return server.NewHandlerError(response.StatusIntervalServerError, []byte("My mistake\n"))
	default:
		w.Write([]byte("All good\n"))
	}
	return nil
} 
