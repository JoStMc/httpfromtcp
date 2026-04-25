package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/JoStMc/httpfromtcp/internal/headers"
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

func respond200() []byte {
    return []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
} 

func respond400() []byte {
    return []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>your request was bad</p>
  </body>
</html>`)
} 

func respond500() []byte {
    return []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>This one is on me.</p>
  </body>
</html>`)
} 

func handlerPaths(w *response.Writer, req *request.Request) {
	target := req.RequestLine.RequestTarget
	if strings.HasPrefix(target, "/httpbin/") {
		path := strings.Trim(target, "/httpbin/")
		proxy := fmt.Sprintf("https://httpbin.org/%s", path)
		h := headers.NewHeaders()
		h.Set("Connection", "close")
		h.Set("Content-Type", "text/plain")
		h.Set("Transfer-Encoding", "chunked")

		res, err := http.Get(proxy)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteStatusLine(response.StatusOK)
		w.WriteHeaders(h)

		for {
			p := make([]byte, 1024)
			n, err := res.Body.Read(p)
			if err != nil {
				if errors.Is(err, io.EOF) {
					n, _ := w.WriteChunkedBodyDone()
					fmt.Println("Bytes parsed:", n)
					return
				} else {
					log.Fatal(err)
				} 
			}
			_, err = w.WriteChunkedBody(p[:n])
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Bytes parsed: ", n)
		} 
	} 

	var b []byte
	switch target {
	case "/yourproblem":
		b = respond400()
		w.WriteStatusLine(response.StatusBadRequest)
	case "/myproblem":
		b = respond500()
		w.WriteStatusLine(response.StatusIntervalServerError)
	default:
		b = respond200()
		w.WriteStatusLine(response.StatusOK)
	}

	headers := response.GetDefaultHeaders(len(b))
	headers.Replace("Content-Type", "text/html")

	w.WriteHeaders(headers)
	w.WriteBody(b)
} 
