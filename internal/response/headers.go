package response

import (
	"fmt"
	"io"

	"github.com/jondatkins/http_from_tcp/internal/headers"
)

type Writer struct {
	Headers    headers.Headers
	StatusCode StatusCode
	Body       []byte
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	w.StatusCode = statusCode
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	w.Headers = headers
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	w.Body = p
	return len(p), nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Content-Type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}
