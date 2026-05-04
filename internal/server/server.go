package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	portString := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", portString)
	if err != nil {
		return nil, err
	}
	server := &Server{
		listener: listener,
	}
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error acception connection: %v", err)
			continue
		}
		go s.handle(conn)
		// go func(c net.Conn) {
		// 	s.handle(c)
		// }(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	testResp := `HTTP/1.1 200 OK
	Content-Type: text/plain
	Content-Length: 13

	Hello World!`
	conn.Write([]byte(testResp))
	return
}
