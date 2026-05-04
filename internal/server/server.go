package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	Listener net.Listener
	IsClosed atomic.Bool
}

func Serve(port int) (*Server, error) {
	portString := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", portString)
	if err != nil {
		return nil, err
	}
	server := Server{
		Listener: l,
	}
	server.listen()
	return &server, nil
}

func (s *Server) Close() error {
	s.IsClosed.Store(true)
	s.Listener.Close()
	return nil
}

func (s *Server) listen() {
	for {
		// Wait for a connection.
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.IsClosed.Load() {
				break
			}
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			// Echo all incoming data.
			// io.Copy(c, c)
			s.handle(c)
			// Shut down the connection.
			c.Close()
		}(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	testResp := `HTTP/1.1 200 OK
	Content-Type: text/plain
	Content-Length: 13

	Hello World!
	defer conn.Close()`
	fmt.Println("New connection accepted")

	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	conn.Write([]byte(testResp))
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}
	fmt.Println("Received data:", string(buffer))
}
