package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := make([]byte, 8)

	for {
		bytesRead, err := reader.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Println("Error reading file:", err)
			return
		}
		file_line := string(buffer[:bytesRead])
		// trimmed_line := strings.TrimSpace(file_line)
		fmt.Print("read: " + file_line + "\n")
	}
}
