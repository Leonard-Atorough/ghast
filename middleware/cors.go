package middleware

import (
	"strconv"

	"github.com/Leonard-Atorough/ghast"
)

const defaultAllowedMethods = "GET, POST, PUT, DELETE, OPTIONS"

type CorsOptions struct {
	AllowedOrigins    []string
	AllowedMethods    []string
	AllowedHeaders    []string
	PreflightMaxAge   int  // Optional: Max age for preflight requests in seconds
	PreflightContinue bool // Optional: Whether to continue processing preflight requests (default: false)
	Credentials       bool // Optional: Whether to allow credentials (default: false)
}

// CorsMiddleware returns a middleware function that adds CORS headers to responses. It allows all origins by default, but can be configured with specific allowed origins, methods, and headers.
func CorsMiddleware(options CorsOptions) ghast.Middleware {
	return func(next ghast.Handler) ghast.Handler {
		return ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
			originKey, originValue := configureOrigins(options)
			methodKey, methodValue := configureMethods(options)
			headers := configureAllowedHeaders(options, r)
			pfKey, pfValue := configurePreflightMaxAge(options)
			credKey, credValue := configureCredentials(options)

			w.SetHeader(originKey, originValue)
			w.SetHeader(methodKey, methodValue)
			w.SetHeader(pfKey, pfValue)
			w.SetHeader(credKey, credValue)
			for key, value := range headers {
				w.SetHeader(key, value)
			}
			// Handle preflight requests
			if r.Method == ghast.OPTIONS && r.Headers["Access-Control-Request-Method"] != "" {

				// return 200 OK for preflight unless told to continue processing
				if !options.PreflightContinue {
					w.Status(200)
					w.Send([]byte("OK"))
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func configurePreflightMaxAge(options CorsOptions) (key string, value string) {
	if options.PreflightMaxAge > 0 {
		return "Access-Control-Max-Age", strconv.Itoa(options.PreflightMaxAge)
	}
	return "", ""
}

func configureCredentials(options CorsOptions) (key string, value string) {
	if options.Credentials {
		return "Access-Control-Allow-Credentials", "true"
	}
	return "", ""
}

func configureOrigins(options CorsOptions) (key string, value string) {
	origins := options.AllowedOrigins

	if len(origins) == 0 {
		return "Access-Control-Allow-Origin", "*"
	}
	originsStr := ""
	for _, origin := range origins {
		originsStr += origin + ", "
	}
	// Remove trailing comma and space
	if len(originsStr) > 2 {
		originsStr = originsStr[:len(originsStr)-2]
	}
	return "Access-Control-Allow-Origin", originsStr
}

func configureMethods(options CorsOptions) (key string, value string) {
	methods := options.AllowedMethods

	if len(methods) == 0 {
		return "Access-Control-Allow-Methods", defaultAllowedMethods
	}
	methodsStr := ""
	for _, method := range methods {
		methodsStr += method + ", "
	}
	// Remove trailing comma and space
	if len(methodsStr) > 2 {
		methodsStr = methodsStr[:len(methodsStr)-2]
	}
	return "Access-Control-Allow-Methods", methodsStr
}

func configureAllowedHeaders(options CorsOptions, r *ghast.Request) map[string]string {
	allowedHeadersStr := ""
	allowedHeaders := options.AllowedHeaders
	headers := make(map[string]string)

	if len(allowedHeaders) == 0 {
		allowedHeadersStr = r.Headers["Access-Control-Request-Headers"]
		headers["Vary"] = "Access-Control-Request-Headers"
	} else {
		allowedHeadersStr := ""
		for _, header := range allowedHeaders {
			allowedHeadersStr += header + ", "
		}
		// Remove trailing comma and space
		if len(allowedHeadersStr) > 2 {
			allowedHeadersStr = allowedHeadersStr[:len(allowedHeadersStr)-2]
		}
	}
	headers["Access-Control-Allow-Headers"] = allowedHeadersStr

	return headers
}
