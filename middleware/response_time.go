package middleware

import (
	"fmt"
	ghast "ghast/lib"
	"time"
)

const defaultResponseTimeHeader = "X-Response-Time"
const defaultResponseTimeSuffix = "ms"

// ResponseTimeOptions defines the options for the ResponseTimeMiddleware.
type ResponseTimeOptions struct {
	HeaderName string // The name of the header to set the response time in (default: "X-Response-Time")
	Suffix     string // Optional suffix to append to the response time (e.g. "ms" for milliseconds)
}

// ResponseTimeMiddleware is a middleware that measures the time taken to process a request and sets it in the response header.
//
// Options:
//   - HeaderName: The name of the header to set the response time in (default: "X-Response-Time")
//   - Suffix: Optional suffix to append to the response time (e.g. "ms" for milliseconds)
func ResponseTimeMiddleware(opts ResponseTimeOptions) ghast.Middleware {
	headerName := defaultResponseTimeHeader
	if opts.HeaderName != "" {
		headerName = opts.HeaderName
	}

	// time.Since returns a Duration, which is in nanoseconds. We can convert it to the desired unit based on the suffix. For example, if suffix is "ms", we divide by time.Millisecond to get milliseconds.
	modifier := getModifierForSuffix(opts.Suffix)

	suffix := defaultResponseTimeSuffix
	if opts.Suffix != "" {
		suffix = opts.Suffix
	}
	return func(next ghast.Handler) ghast.Handler {
		return ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			w.SetHeader(headerName, fmt.Sprintf("%d%s", duration/modifier, suffix))
		})

	}
}

func getModifierForSuffix(suffix string) time.Duration {
	switch suffix {
	case "ms":
		return time.Millisecond
	case "s":
		return time.Second
	case "us":
		return time.Microsecond
	case "ns":
		return time.Nanosecond
	default:
		return time.Millisecond // Default to milliseconds if suffix is unrecognized
	}
}
