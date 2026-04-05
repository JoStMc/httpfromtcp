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
	defer f.Close()

	var currentLine string
	for {
		b := make([]byte, 8)
		n, err := f.Read(b)
		if err != nil {
			if errors.Is(err, io.EOF){
				break
			} 
			fmt.Printf("error: %s\n", err.Error())
			break
		}

		str := string(b[:n])
		lines := strings.Split(str, "\n")
		currentLine = currentLine + lines[0]
		if len(lines) > 1 {
			for i := 0; i < len(lines) - 1; i++ {
				fmt.Printf("read: %s\n", currentLine)
				currentLine = lines[i+1]
			} 
		} 
	} 
}
