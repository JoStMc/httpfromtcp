package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

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
		fmt.Printf("read: %s\n", str)
	} 
}
