package ghast

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

// ResponseWriter interface for constructing and sending HTTP responses.
type ResponseWriter interface {
	Header() map[string]string // Returns the response headers map for setting headers before writing the body.

	Status(statusCode int) ResponseWriter // Sets the HTTP status code and returns self for chaining.

	SetHeader(key, value string) ResponseWriter // SetHeader sets a response header and returns self for chaining.

	Send([]byte) (int, error) // Send writes the data with optional content-type detection.

	SendString(string) (int, error) // SendString writes a string response.

	JSON(statusCode int, data interface{}) error // JSON marshals data as JSON and sends it with application/json content-type.

	JSONPretty(statusCode int, data interface{}) error // JSONPretty marshals data as pretty-printed JSON.
}

// responseWriter implements ResponseWriter interface.
type responseWriter struct {
	conn       net.Conn
	headers    map[string]string
	statusCode int
	statusText string
	written    bool // Tracks whether status/headers have been written
}

// NewResponseWriter creates a new ResponseWriter for the given connection.
func NewResponseWriter(conn net.Conn) ResponseWriter {
	return &responseWriter{
		conn:       conn,
		headers:    make(map[string]string),
		statusCode: 200,
		statusText: "OK",
		written:    false,
	}
}

// Header returns the response headers map.
func (rw *responseWriter) Header() map[string]string {
	return rw.headers
}

// Status sets the HTTP status code and returns self for chaining.
func (rw *responseWriter) Status(statusCode int) ResponseWriter {
	if !rw.written {
		rw.statusCode = statusCode
		rw.statusText = httpStatusText(statusCode)
	}
	return rw
}

// SetHeader sets a response header and returns self for chaining.
func (rw *responseWriter) SetHeader(key, value string) ResponseWriter {
	rw.headers[key] = value
	return rw
}

// WriteHeader sets the HTTP status code (non-chainable, called automatically by Write).
// @internal This is not meant to be called directly by handlers. Use Status() for chaining instead.
func (rw *responseWriter) writeHeader(statusCode int) {
	if !rw.written {
		rw.statusCode = statusCode
		rw.statusText = httpStatusText(statusCode)
	}
}

// Write writes data to the response body.
// @internal - This is called by Send() and SendString() to write the response body. It automatically writes the status line and headers if they haven't been written yet.
func (rw *responseWriter) write(data []byte) (int, error) {
	if !rw.written {
		rw.writeStatusAndHeaders()
		rw.written = true
	}
	return rw.conn.Write(data)
}

// Send writes data to the response body.
func (rw *responseWriter) Send(data []byte) (int, error) {
	return rw.write(data)
}

// SendString writes a string response.
func (rw *responseWriter) SendString(s string) (int, error) {
	return rw.write([]byte(s))
}

// JSON marshals data as JSON and sends it with application/json content-type.
func (rw *responseWriter) JSON(statusCode int, data interface{}) error {
	rw.Status(statusCode)
	rw.SetHeader("Content-Type", "application/json")

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = rw.write(jsonData)
	return err
}

// JSONPretty marshals data as pretty-printed JSON.
func (rw *responseWriter) JSONPretty(statusCode int, data interface{}) error {
	rw.Status(statusCode)
	rw.SetHeader("Content-Type", "application/json")

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	_, err = rw.write(jsonData)
	return err
}

// writeStatusAndHeaders writes the HTTP status line and headers.
func (rw *responseWriter) writeStatusAndHeaders() {
	var buf strings.Builder
	fmt.Fprintf(&buf, "HTTP/1.1 %d %s\r\n", rw.statusCode, rw.statusText)
	for key, value := range rw.headers {
		fmt.Fprintf(&buf, "%s: %s\r\n", key, value)
	}
	buf.WriteString("\r\n")
	rw.conn.Write([]byte(buf.String()))
}
