package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const inputFilePath = "messages.txt"

func main() {
	// fmt.Printf("Reading data from %s\n", inputFilePath)
	// fmt.Println("=============================================")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-interrupt
		fmt.Println("Quitting")
		os.Exit(0)
	}()
	fmt.Println("Program running, press ctrl+c to exit")
	ln, err := net.Listen("tcp", ":42069")
	defer ln.Close()
	if err != nil {
		log.Fatalf("could not create listener", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("error accepting connection", err)
		}
		fmt.Println("Connection Accepted", conn)

		listenerContents := getLinesChannel(conn)
		for line := range listenerContents {
			fmt.Println("read:", line)
		}
	}
	// fmt.Println("Connection closed")
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)
		defer f.Close()
		currentLine := ""
		buffer := make([]byte, 8, 8)
		for {
			bytesRead, err := f.Read(buffer)
			if err != nil {
				if currentLine != "" {
					lines <- currentLine
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}
			str := string(buffer[:bytesRead])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- currentLine + parts[i]
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}
	}()
	return lines
}
