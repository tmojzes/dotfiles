// Command armada-mcp is an MCP server that runs readonly kubectl/oc commands
// via a Jenkins job.
//
// Ported from the former Python implementation
// (opencode/.config/opencode/mcp/armada.py).
//
// Environment contract (registered in the global opencode.json under
// mcp.servers.armada.environment; hand-exported from ~/keys):
//
//	JENKINS_USER     Jenkins user for the armada job
//	JENKINS_API_KEY  Jenkins API token (~/keys/jenkins-api-key)
package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// serverConfig holds the polling and HTTP timing knobs (CLI flags).
type serverConfig struct {
	queueInterval  time.Duration
	buildInterval  time.Duration
	requestTimeout time.Duration
}

func main() {
	var cfg serverConfig
	flag.DurationVar(&cfg.queueInterval, "queue-interval", 5*time.Second, "time between Jenkins queue polls")
	flag.DurationVar(&cfg.buildInterval, "build-interval", 10*time.Second, "time between Jenkins build polls")
	flag.DurationVar(&cfg.requestTimeout, "request-timeout", 30*time.Second, "HTTP timeout for Jenkins requests")
	flag.Parse()

	if err := serve(context.Background(), os.Stdin, os.Stdout, cfg); err != nil {
		fmt.Fprintln(os.Stderr, "armada-mcp:", err)
		os.Exit(1)
	}
}

// jsonSyntaxError marks input that failed JSON decoding; serve answers it
// with a -32700 response instead of treating it as fatal.
type jsonSyntaxError struct{ msg string }

func (e *jsonSyntaxError) Error() string { return e.msg }

type messageReader struct{ r *bufio.Reader }

// readMessage reads one JSON-RPC message and reports the framing it was
// delivered with ("newline" or "content-length"). Blank lines are skipped;
// io.EOF signals a closed input stream.
func (m *messageReader) readMessage() (string, any, error) {
	for {
		line, readErr := m.r.ReadBytes('\n')
		if readErr != nil && readErr != io.EOF {
			return "", nil, readErr
		}
		trimmed := bytes.TrimRight(line, "\r\n")
		if len(trimmed) == 0 {
			if readErr == io.EOF {
				return "", nil, io.EOF
			}
			continue
		}

		if strings.HasPrefix(string(trimmed), "Content-Length:") {
			return m.readFramedMessage(trimmed)
		}

		msg, err := decodeJSON(trimmed)
		if err != nil {
			return "newline", nil, err
		}
		return "newline", msg, nil
	}
}

func (m *messageReader) readFramedMessage(first []byte) (string, any, error) {
	headers := map[string]string{}
	if key, value, ok := splitHeaderLine(first); ok {
		headers[key] = value
	}
	for {
		line, readErr := m.r.ReadBytes('\n')
		if readErr != nil && readErr != io.EOF {
			return "", nil, readErr
		}
		trimmed := bytes.TrimRight(line, "\r\n")
		if len(trimmed) == 0 {
			break
		}
		if key, value, ok := splitHeaderLine(trimmed); ok {
			headers[key] = value
		}
		if readErr == io.EOF {
			return "", nil, io.EOF
		}
	}
	contentLength, err := strconv.Atoi(headers["content-length"])
	if err != nil {
		return "", nil, errors.New("invalid Content-Length header")
	}
	body := make([]byte, contentLength)
	if _, err := io.ReadFull(m.r, body); err != nil {
		return "", nil, io.EOF
	}
	msg, err := decodeJSON(body)
	if err != nil {
		return "content-length", nil, err
	}
	return "content-length", msg, nil
}

func splitHeaderLine(line []byte) (string, string, bool) {
	key, value, found := strings.Cut(string(line), ":")
	if !found {
		return "", "", false
	}
	return strings.ToLower(strings.TrimSpace(key)), strings.TrimSpace(value), true
}

func decodeJSON(data []byte) (any, error) {
	var msg any
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, &jsonSyntaxError{msg: err.Error()}
	}
	return msg, nil
}

type messageWriter struct{ w *bufio.Writer }

func (mw *messageWriter) write(payload any, framing string) error {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if framing == "content-length" {
		if _, err := fmt.Fprintf(mw.w, "Content-Length: %d\r\n\r\n", len(encoded)); err != nil {
			return err
		}
	} else {
		encoded = append(encoded, '\n')
	}
	if _, err := mw.w.Write(encoded); err != nil {
		return err
	}
	return mw.w.Flush()
}

// serve runs the stdio JSON-RPC loop: read a message (or batch), answer every
// request that carries an id, and mirror the inbound framing on responses.
func serve(ctx context.Context, in io.Reader, out io.Writer, cfg serverConfig) error {
	reader := &messageReader{r: bufio.NewReader(in)}
	writer := &messageWriter{w: bufio.NewWriter(out)}

	for {
		framing, message, err := reader.readMessage()
		if err != nil {
			var syntaxErr *jsonSyntaxError
			if errors.As(err, &syntaxErr) {
				resp := errorResponse(nil, &rpcError{Code: -32700, Message: "Invalid JSON: " + syntaxErr.msg})
				if werr := writer.write(resp, framing); werr != nil {
					return werr
				}
				continue
			}
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("fatal reader error: %w", err)
		}

		messages, isBatch := message.([]any)
		if !isBatch {
			messages = []any{message}
		}

		var responses []any
		for _, item := range messages {
			request, ok := item.(map[string]any)
			if !ok {
				responses = append(responses, errorResponse(nil, &rpcError{Code: -32600, Message: "Invalid Request"}))
				continue
			}

			method, methodIsString := request["method"].(string)
			id := request["id"]
			if !methodIsString {
				if id != nil {
					responses = append(responses, errorResponse(id, &rpcError{Code: -32600, Message: "Invalid Request"}))
				}
				continue
			}

			result, rpcErr := handleRequest(ctx, cfg, method, request["params"])
			if rpcErr != nil {
				if id != nil {
					responses = append(responses, errorResponse(id, rpcErr))
				}
				continue
			}
			if id == nil || result == nil {
				continue
			}
			responses = append(responses, successResponse(id, result))
		}

		if len(responses) == 0 {
			continue
		}
		payload := any(responses[0])
		if len(responses) > 1 {
			payload = responses
		}
		if err := writer.write(payload, framing); err != nil {
			return err
		}
	}
}
