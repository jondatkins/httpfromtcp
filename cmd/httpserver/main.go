// package httpserver
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/jondatkins/http_from_tcp/internal/headers"
	"github.com/jondatkins/http_from_tcp/internal/request"
	"github.com/jondatkins/http_from_tcp/internal/response"
	"github.com/jondatkins/http_from_tcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/" {
		body := []byte(`
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
		`)
		headers := response.GetDefaultHeaders(len(body))
		headers.Replace("Content-Type", "text/html")
		w.WriteStatusLine(response.StatusCodeSuccess)
		w.WriteHeaders(headers)
		w.WriteBody([]byte(body))
	}
	if req.RequestLine.RequestTarget == "/yourproblem" {
		body := []byte(`
<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
		`)
		headers := response.GetDefaultHeaders(len(body))
		headers.Replace("Content-Type", "text/html")
		w.WriteStatusLine(response.StatusCodeBadRequest)
		w.WriteHeaders(headers)
		w.WriteBody([]byte(body))
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		body := []byte(`
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
		`)
		headers := response.GetDefaultHeaders(len(body))
		headers.Replace("Content-Type", "text/html")
		w.WriteStatusLine(response.StatusCodeInternalServerError)
		w.WriteHeaders(headers)
		w.WriteBody([]byte(body))
	}

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		path := strings.TrimPrefix(
			req.RequestLine.RequestTarget, "/httpbin",
		)

		url := "https://httpbin.org/" + path

		resp, err := http.Get(url)
		if err != nil {
			w.WriteStatusLine(response.StatusCodeInternalServerError)
			return
		}
		defer resp.Body.Close()

		respHeaders := response.GetDefaultHeaders(0)
		respHeaders.Set("Trailer", "X-Content-SHA256")
		respHeaders.Set("Trailer", "X-Content-Length")
		respHeaders.Replace("Content-Type", "text/plain")
		respHeaders.Replace("Transfer-Encoding", "chunked")
		respHeaders.Delete("Content-Length")

		w.WriteStatusLine(response.StatusCodeSuccess)
		w.WriteHeaders(respHeaders)

		buf := make([]byte, 1024)
		var fullBody bytes.Buffer

		for {
			n, err := resp.Body.Read(buf)

			fmt.Println(n)

			if n > 0 {
				chunk := buf[:n]
				fullBody.Write(chunk)

				_, writeErr := w.WriteChunkedBody(chunk)
				if writeErr != nil {
					return
				}
				// _, writeErr := w.WriteChunkedBody(buf[:n])
				// if writeErr != nil {
				// 	return
				// }
			}

			if err == io.EOF {
				break
			}

			if err != nil {
				return
			}
		}
		hash := sha256.Sum256(fullBody.Bytes())

		trailers := make(headers.Headers)

		trailers.Set(
			"X-Content-SHA256",
			hex.EncodeToString(hash[:]),
		)

		trailers.Set("X-Content-Length", strconv.Itoa(fullBody.Len()))

		w.WriteChunkedBodyDone(trailers)
	}
}
