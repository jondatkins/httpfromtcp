package server

import (
	"bytes"
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
	Handler func(w io.Writer, req *request.Request) *HandlerError
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

func (he HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, he.StatusCode)
	messageBytes := []byte(he.Message)
	headers := response.GetDefaultHeaders(len(messageBytes))
	response.WriteHeaders(w, headers)
	w.Write(messageBytes)
}

// func writeError(writer io.Writer, handlerError HandlerError) error {
// 	// _, err := writer.Write([]byte(fmt.Sprintf(handlerError.statusCode + " " + handlerError.message)))
// 	body := handlerError.Message
//
// 	err := response.WriteStatusLine(writer, handlerError.StatusCode)
// 	if err != nil {
// 		return err
// 	}
//
// 	headers := response.GetDefaultHeaders(len(body))
//
// 	err = response.WriteHeaders(writer, headers)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = writer.Write([]byte(body))
// 	return err
// }

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

	req, err := request.RequestFromReader(conn)
	if err != nil {

		// Write(conn, HandlerError{
		hErr := &HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message:    "Bad Request\n",
		}
		hErr.Write(conn)
		return
	}

	// var buffer bytes.Buffer
	buffer := bytes.NewBuffer([]byte{})

	handlerErr := s.handler(buffer, req)
	if handlerErr != nil {
		// writeError(conn, *handlerErr)
		handlerErr.Write(conn)
		return
	}
	b := buffer.Bytes()
	response.WriteStatusLine(conn, response.StatusCodeSuccess)
	headers := response.GetDefaultHeaders(len(b))
	response.WriteHeaders(conn, headers)
	conn.Write(b)
	return
}
