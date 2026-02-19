package ghast

// Handler interface defines how HTTP requests are processed.
// Implementations should read from the Request and write to the ResponseWriter.
// For those familiar with .NET, this is similar to the IHttpHandler interface.
type Handler interface {
	// ServeHTTP processes an HTTP request and constructs an appropriate response.
	// The ResponseWriter is used to write the response back to the client,
	// while the Request contains all the information about the incoming HTTP request.
	ServeHTTP(ResponseWriter, *Request)
}

// HandlerFunc is a type that allows ordinary functions to be used as HTTP handlers.
// It implements the Handler interface by defining a ServeHTTP method that calls the underlying function.
//
// Example:
//   handler := ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
//       w.JSON(200, map[string]string{"message": "Hello"})
//   })
type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP implements the Handler interface for HandlerFunc.
func (f HandlerFunc) ServeHTTP(rw ResponseWriter, req *Request) {
	f(rw, req)
}
