package ghast

// HTTPError represents an HTTP error with status code and message.
type HTTPError struct {
	StatusCode int    `json:"status"`
	Message    string `json:"error"`
}

// Error sends an error response as JSON with the given status code and message.
func Error(rw ResponseWriter, statusCode int, message string) error {
	errResp := HTTPError{
		StatusCode: statusCode,
		Message:    message,
	}
	return rw.JSON(statusCode, errResp)
}

// ErrorString sends an error response with a string body.
func ErrorString(rw ResponseWriter, statusCode int, message string) error {
	rw.Status(statusCode).SetHeader("Content-Type", "text/plain")
	_, err := rw.SendString(message)
	return err
}
