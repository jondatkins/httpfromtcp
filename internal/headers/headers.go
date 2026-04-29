package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n\r\n"

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	numBytes := len(data)
	header := strings.TrimRight(string(data), " ")
	// Look for a CRLF, if it doesn't find one, assume you haven't been
	// given enough data yet. Consume no data, return false for done, and nil for err.
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	// If you do find a CRLF, but it's at the start of the data, you've found the end of the headers, so return the proper values immediately.
	if idx == 0 {
		// numBytes / 0 ?
		return numBytes, true, nil
	}
	// Host: localhost:42069\r\n\r\n
	// Remove any extra whitespace from the key and value, but ensure there are no spaces between the colon and the key.

	headerParts := strings.Split(header, " ")
	if len(headerParts) < 2 {
		return 0, false, fmt.Errorf("header line should have at least 2 parts: %s", header)
	}
	if headerParts[0] == "" {
		return 0, false, fmt.Errorf("field name must not be preceded by whitepace: %s", header)
	}
	key := headerParts[0]
	if key[len(key)-1:] != ":" {
		return 0, false, fmt.Errorf("field name must end in ':'")
	}
	// Assuming the format was valid (if it isn't return an error), add the key/value pair to the Headers map and return the number of bytes consumed.
	value := headerParts[1]
	key = key[:len(key)-1]
	crlfLength := len(crlf)
	value = value[:len(value)-crlfLength]
	h[key] = value
	// Note: The Parse function should only return done=true when the data starts with a CRLF, which can't happen when it finds a new key/value pair.
	return numBytes - crlfLength/2, false, nil
}

func headerLineFromString(str string) (string, error) {
	parts := strings.Split(str, " ")
	if len(parts) < 2 {
		return "", fmt.Errorf("poorly formatted header-line: %s", str)
	}

	key := parts[0]
	// for _, c := range key {
	// 	if c < 'A' || c > 'Z' {
	// 		return "", fmt.Errorf("invalid method: %s", key)
	// 	}
	// }

	value := parts[1]
	return key + " " + value, nil
}
