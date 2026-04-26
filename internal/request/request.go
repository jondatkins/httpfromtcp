package request

import (
	"fmt"
	"io"
	"log"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const httpVersion = "1.1"

func RequestFromReader(reader io.Reader) (*Request, error) {
	req, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalf("%sn", err)
	}
	reqLineList, err := parseRequestLine(string(req))
	if err != nil {
		return &Request{}, fmt.Errorf("Error parsing request line: %s", string(req))
	}
	request := Request{
		RequestLine: RequestLine{
			HttpVersion:   reqLineList[2],
			RequestTarget: reqLineList[1],
			Method:        reqLineList[0],
		},
	}
	return &request, nil
}

func parseRequestLine(line string) ([]string, error) {
	requestLines := strings.Split(line, "\r\n")
	if len(requestLines) == 0 {
		return nil, fmt.Errorf("No new line in string: %s", line)
	}
	requestLine := requestLines[0]
	// Should be e.g. 'POST /coffee HTTP/1.1'
	// GET / HTTP/1.1
	reqLineParts := strings.Split(requestLine, " ")
	if len(reqLineParts) < 3 {
		return nil, fmt.Errorf("Request line should have 3 parts: %s", reqLineParts)
	}
	// check 3rd string is only alphabetic
	isAlphabetic := true
	for _, char := range reqLineParts[0] {
		if !unicode.IsLetter(char) {
			isAlphabetic = false
			break
		}
	}
	if !isAlphabetic {
		return nil, fmt.Errorf("Method name must contain letters only: %s", reqLineParts[2])
	}

	if reqLineParts[0] != strings.ToUpper(reqLineParts[0]) {
		return nil, fmt.Errorf("Method name must be all caps: %s", reqLineParts[0])
	}
	httpVersionList := strings.Split(reqLineParts[2], "/")
	if len(httpVersionList) < 2 {
		return nil, fmt.Errorf("HTTP version string should contain '/': %s", reqLineParts[2])
	}
	if httpVersionList[1] != httpVersion {
		return nil, fmt.Errorf("HTTP version should be: %s", httpVersion)
	}
	reqLineParts[2] = httpVersionList[1]
	fmt.Println(reqLineParts)
	return reqLineParts, nil
}
