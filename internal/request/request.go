package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type parserState int

const (
	initialized parserState = iota
	done
)

type Request struct {
	RequestLine RequestLine
	parserState parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	crlf       = "\r\n"
	bufferSize = 8
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, bufferSize)
	readToIndex := 0
	// parsedFromIndex := 0
	request := Request{
		RequestLine: RequestLine{
			HttpVersion:   "",
			RequestTarget: "",
			Method:        "",
		},
		parserState: initialized,
	}
	for request.parserState != done {
		if readToIndex == cap(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}
		n, err := reader.Read(buffer[readToIndex:])
		if err != nil {
			if err == io.EOF {
				if readToIndex == 0 {
					break
				}
				request.parserState = done
				break
			}
			return &Request{}, err
		}
		readToIndex += n
		consumed, err := request.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}

		if consumed > readToIndex {
			return nil, fmt.Errorf("parser consumed more bytes than available")
		}

		if consumed > 0 {
			copy(buffer, buffer[consumed:readToIndex])
			readToIndex -= consumed
		}

	}
	return &request, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		// return nil, fmt.Errorf("could not find CRLF in request-line")
		return &RequestLine{}, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return &RequestLine{}, 0, err
	}
	return requestLine, idx + len(crlf), nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")

	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.parserState {
	case initialized:
		requestLine, bytesParsed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if bytesParsed == 0 {
			return 0, nil
		}
		if bytesParsed > 0 {
			// r.RequestLine.HttpVersion = requestLine.HttpVersion
			// r.RequestLine.RequestTarget = requestLine.RequestTarget
			// r.RequestLine.Method = requestLine.Method
			r.RequestLine = *requestLine
			r.parserState = done
		}
		return bytesParsed, nil
	case done:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}
