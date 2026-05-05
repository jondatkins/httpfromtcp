package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jondatkins/http_from_tcp/internal/headers"
)

type StatusCode int

const (
	OK         StatusCode = 200
	BadReq     StatusCode = 400
	IntServErr StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var line string

	switch statusCode {
	case OK:
		line = "HTTP/1.1 200 OK\r\n"
	case BadReq:
		line = "HTTP/1.1 400 Bad Request\r\n"
	case IntServErr:
		line = "HTTP/1.1 500 Internal Server Error\r\n"
	default:
		line = "HTTP/1.1 "
	}
	_, err := io.WriteString(w, line)
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	defaultHeaders := headers.NewHeaders()
	defaultHeaders["Content-Length"] = strconv.Itoa(contentLen)
	defaultHeaders["Connection"] = "close"
	defaultHeaders["Content-Type"] = "text/plain"
	return defaultHeaders
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := io.WriteString(w, key+": "+value+"\r\n")
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
