package ghast

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

const CRLF = "\r\n"
const HTTPVersion = "HTTP/1.1"

// HTTP Methods
const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
)

// Request represents an HTTP request with parsed components.
type Request struct {
	Method  string            // HTTP method (GET, POST, etc.)
	Path    string            // URL path (without query string)
	Headers map[string]string // HTTP headers
	Body    string            // Request body as string
	Version string            // HTTP version (e.g., "HTTP/1.1")
	Params  map[string]string // Route parameters (e.g., from path variables)
	Queries map[string]string // Query parameters
	ClientIP string            // Client IP address (to be populated by server)
}

// Query retrieves a query parameter by key. Returns empty string if not found.
func (r *Request) Query(key string) string {
	if r.Queries == nil {
		return ""
	}
	return r.Queries[key]
}

// Param retrieves a route parameter by key. Returns empty string if not found.
func (r *Request) Param(key string) string {
	if r.Params == nil {
		return ""
	}
	return r.Params[key]
}

// JSON unmarshals the request body as JSON into the provided value. To be replaced by a common body command and content type handling in the future.
func (r *Request) JSON(v any) error {
	return json.Unmarshal([]byte(r.Body), v)
}

// GetHeader retrieves a header value (case-insensitive).
// Returns empty string if header not found.
func (r *Request) GetHeader(key string) string {
	// Try exact match first
	if val, ok := r.Headers[key]; ok {
		return val
	}
	// Try case-insensitive match
	for k, v := range r.Headers {
		if strings.EqualFold(k, key) {
			return v
		}
	}
	return ""
}

// ContentType returns the Content-Type header value.
func (r *Request) ContentType() string {
	return r.GetHeader("Content-Type")
}

// ParseRequest parses a raw HTTP request string into a Request struct.
// It extracts the method, path, version, headers, and query parameters.
//
// Note: Body parsing is handled separately by the server due to TCP/stream considerations.
//
// TODO:
//   - Add support for different content types and encodings in the request body.
//   - Add support for duplicate headers and query parameters.
//   - Add validation for header names and values.
//   - Add support for URL decoding of query parameters.
func ParseRequest(rawRequest string) (*Request, error) {
	lines := strings.Split(rawRequest, CRLF)

	if len(lines) < 1 {
		return nil, fmt.Errorf("invalid request: no lines found")
	}

	method, path, version, err := parseRequestLine(lines)
	if err != nil {
		return nil, err
	}

	headers, err := parseHeaders(lines[1:])
	if err != nil {
		return nil, err
	}

	var queries map[string]string
	if strings.Contains(path, "?") {
		var err error
		queries, err = parseParams(strings.Split(path, "?")[1])
		if err != nil {
			return nil, err
		}
		path = strings.Split(path, "?")[0] // Strip query string from path for routing
	}

	var params map[string]string // Params will be populated later by the router when matching dynamic routes

	req := &Request{
		Method:  method,
		Path:    path,
		Version: version,
		Headers: headers,
		Queries: queries,
		Params:  params,
		Body:    "", // Body will be populated later if Content-Length is present
	}
	return req, nil
}

// parseRequestLine parses an HTTP request line (e.g., "GET /index.html HTTP/1.1") into method, path, and version.
//
// TODO:
//   - Add validation to ensure the method and version are valid.
//   - Add support for parsing the path to separate it from query parameters.
func parseRequestLine(lines []string) (method, path, version string, err error) {
	requestLine := lines[0]
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid request line: expected 3 parts, got %d", len(parts))
	}
	method = parts[0]
	path = parts[1]
	version = parts[2]

	validMethods := []string{GET, POST, PUT, DELETE, HEAD, OPTIONS, PATCH}

	isValidMethod := slices.Contains(validMethods, strings.ToUpper(method))

	if !isValidMethod {
		return "", "", "", fmt.Errorf("invalid request line: unknown method %s", method)
	}

	return method, path, version, nil
}

// parseHeaders parses HTTP request header lines into a map of header names to values.
//
// TODO:
//   - Add support for handling duplicate headers.
//   - Add validation for header names and values.
//   - Add support for handling different line endings.
func parseHeaders(lines []string) (map[string]string, error) {
	headers := make(map[string]string)
	for _, line := range lines {
		if line == "" {
			break // End of headers
		}
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header line: %s", line)
		}

		if !isValidHeaderName(parts[0]) || !isValidHeaderValue(parts[1]) {
			return nil, fmt.Errorf("invalid header line: %s", line)
		}
		headers[parts[0]] = parts[1]
	}
	return headers, nil
}

// parseParams parses a query parameter string (e.g., "key1=value1&key2=value2") into a map of key-value pairs.
//
// TODO:
//   - Add support for URL-decoding keys and values.
//   - Add support for handling duplicate query parameter keys.
//   - Add validation for query parameter keys and values.
func parseParams(paramString string) (map[string]string, error) {
	params := make(map[string]string)
	pairs := strings.SplitSeq(paramString, "&")
	for pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid query parameter: %s", pair)
		}
		params[kv[0]] = kv[1]
	}
	return params, nil
}

func isValidHeaderName(name string) bool {
	// Basic validation for header names (can be expanded as needed)
	return name != "" && !strings.ContainsAny(name, " \t\r\n") && !strings.Contains(name, ":") && !strings.HasPrefix(name, " ") && !strings.HasSuffix(name, " ")
}

func isValidHeaderValue(value string) bool {
	// Basic validation for header values (can be expanded as needed)
	return value != "" && !strings.ContainsAny(value, "\r\n")
}
