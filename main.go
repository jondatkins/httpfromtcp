package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

func main() {
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("could not open %s: %s", inputFilePath, err)
	}

	fmt.Printf("Reading data from %s\n", inputFilePath)
	fmt.Println("=============================================")

	fileContents := getLinesChannel(file)
	for line := range fileContents {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		defer f.Close()
		currentLine := ""
		buffer := make([]byte, 8, 8)
		for {
			bytesRead, err := f.Read(buffer)
			if err != nil {
				if currentLine != "" {
					// fmt.Printf("read: %s\n", currentLine)
					ch <- currentLine
					currentLine = ""
				}
				if errors.Is(err, io.EOF) {
					return
				}
				fmt.Printf("error: %s\n", err.Error())
				break
			}
			str := string(buffer[:bytesRead])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				// fmt.Printf("read: %s%s\n", currentLine, parts[i])
				ch <- currentLine + parts[i]
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}
	}()
	return ch
}
