package main

import (
	"fmt"
	"log"
	"net"

	"github.com/JoStMc/httpfromtcp/internal/request"
)

func main() {
	tcpListener, err := net.Listen("tcp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}
	defer tcpListener.Close()


	for {
		connection, err := tcpListener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Connection to message has been accepted on port 42069")

		req, err := request.RequestFromReader(connection)
		if err != nil {
			log.Fatal(err)
		}


		printReqLine(req)
		printHeaders(req)

		fmt.Println("Connection has been closed")
	}
}

func printReqLine(req *request.Request) {
	fmt.Println("Request line:")
	fmt.Println("- Method:", req.RequestLine.Method)
	fmt.Println("- Target:", req.RequestLine.RequestTarget)
	fmt.Println("- Version:", req.RequestLine.HttpVersion)
} 

func printHeaders(req *request.Request) {
	fmt.Println("Headers:")
	for key, value := range req.Headers {
		fmt.Printf("- %s: %s\n", key, value)
	} 
} 
