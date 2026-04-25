package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const port = ":42069"

func main() {
	localAddress, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	remoteAddress, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", localAddress, remoteAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		text, _ := reader.ReadString('\n')
		_, err := conn.Write([]byte(text))
		if err != nil {
			fmt.Printf("Error: %v \n", err)
			continue
		}
	}
}
