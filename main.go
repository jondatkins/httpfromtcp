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
		log.Fatal("could not open %s: %s\n", inputFilePath, err)
	}
	defer file.Close()

	fmt.Printf("Reading data from %s\n", inputFilePath)
	fmt.Println("=============================================")

	// reader := bufio.NewReader(file)
	// buffer := make([]byte, 8)
	currentLine := ""
	for {
		buffer := make([]byte, 8, 8)
		bytesRead, err := file.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("error: %s\n", err.Error())
			break
		}
		byteOfFile := string(buffer[:bytesRead])
		parts := strings.Split(byteOfFile, "\n")
		if len(parts) < 2 {
			currentLine += byteOfFile
		}

		for i := 0; i < len(parts)-1; i++ {
			fmt.Printf("read: %s\n", currentLine+parts[i])
			currentLine = ""
			currentLine += parts[len(parts)-1]
		}
	}
	if len(currentLine) > 0 {
		fmt.Printf("read: %s\n", currentLine)
	}
}
