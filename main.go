package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}

	linesCh := getLinesChannel(f)
	for line := range linesCh {
		fmt.Printf("read: %s\n", line)
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
