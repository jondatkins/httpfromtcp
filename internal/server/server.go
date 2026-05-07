package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/jondatkins/http_from_tcp/internal/request"
	"github.com/jondatkins/http_from_tcp/internal/response"
)

type (
	HandlerError struct {
		StatusCode response.StatusCode
		Message    string
	}
	// Handler func(w io.Writer, req *request.Request) *HandlerError
	Handler func(w *response.Writer, req *request.Request)
)

type Server struct {
	handler  Handler
	listener net.Listener
	closed   atomic.Bool
}

func (he HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, he.StatusCode)
	messageBytes := []byte(he.Message)
	headers := response.GetDefaultHeaders(len(messageBytes))
	response.WriteHeaders(w, headers)
	w.Write(messageBytes)
}

func Serve(port int, handlerFunc Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{
		listener: listener,
		handler:  handlerFunc,
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
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	w := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteStatusLine(response.StatusCodeBadRequest)
		body := []byte(fmt.Sprintf("Error parsing request: %v", err))
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}
	s.handler(w, req)
}
