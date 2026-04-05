package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	tcpListener, err := net.Listen("tcp", "127.0.0.1:42069")
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

		linesCh := getLinesChannel(connection)
		for line := range linesCh {
			fmt.Println(line)
		} 
		fmt.Println("Connection has been closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesCh := make(chan string)

	go func() {
		defer f.Close()
		defer close(linesCh)
		var currentLine string
		for {
			b := make([]byte, 8)
			n, err := f.Read(b)
			if err != nil {
				if errors.Is(err, io.EOF){
					if currentLine != "" {
						linesCh <- currentLine
					} 
					break
				} 
				fmt.Printf("error: %s\n", err.Error())
				break
			}

			str := string(b[:n])
			lines := strings.Split(str, "\n")
			currentLine = currentLine + lines[0]
			for i := 0; i < len(lines) - 1; i++ {
				linesCh <- currentLine
				currentLine = lines[i+1]
			} 
		} 
	}()
	return linesCh
}
