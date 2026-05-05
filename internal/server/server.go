package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/jondatkins/http_from_tcp/internal/response"
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

	err := response.WriteStatusLine(conn, response.OK)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}
	headers := response.GetDefaultHeaders(0)

	err = response.WriteHeaders(conn, headers)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}
	for key, value := range headers {
		fmt.Println(key, " ", value)
	}
	_, err = io.WriteString(conn, "\r\n")
	if err != nil {
		fmt.Printf("error writing final crlf: %v\n", err)
		return
	}
}
