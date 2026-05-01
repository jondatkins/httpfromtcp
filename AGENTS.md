# httpfromtcp

Go project for the [boot.dev](https://boot.dev) HTTP course — builds a raw HTTP/1.1 server and client from TCP/UDP sockets.

## Commands

```
go run ./cmd/tcplistener   # TCP server on :42069
go run ./cmd/udpsender     # UDP sender to :42069
go test ./...              # run all tests
```

Root `main.go` is legacy; ignore it. Real entrypoints are under `cmd/`.

## Architecture

- **cmd/tcplistener** — TCP listener that accepts connections and parses HTTP requests via `internal/request.RequestFromReader`.
- **cmd/udpsender** — stdin-to-UDP client for manual testing.
- **internal/request** — State-machine HTTP request parser reading from an `io.Reader`. Reads small chunks (8 bytes) into a growing buffer. Currently parses the request line only; header parsing is planned per `notes.md`.
- **internal/headers** — Standalone header parser. Header keys are lowercased; duplicate keys are comma-joined. Uses `\r\n` line endings.

## Testing

- Uses `github.com/stretchr/testify`.
- Tests use a custom `chunkReader` to simulate tiny reads (1–50 bytes). This is critical: the parser must handle data arriving in arbitrary chunks.
- Port `42069` is hardcoded everywhere.

## Conventions

- `notes.md` is in `.gitignore` — it's personal scratch notes.
- No linter/formatter configured; follow existing Go conventions (`gofmt`-compatible).
