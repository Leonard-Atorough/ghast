package middleware

import (
	ghast "ghast/lib"
	"time"
)

type RateLimitOptions struct {
	RequestsPerMinute int
}

type rateLimitEntry struct {
	Count     int
	Timestamp int64
}

var rateLimitCollection = make(map[string]rateLimitEntry) // Map of client IP to slice of request timestamps

// RateLimitMiddleware returns a middleware function that implements simple per-IP rate limiting.
// When a new IP is receive, we create a new entry with a timestamp and a coutn. If we receive another request from the same ip, we check if the request timestamp - duration is less than 1 minute. If it is, we check if this plus the count is greater than rpm. If it is, we return a 429 Too Many Requests. if it is not we increment the count. If the request timestamp - duration is greater than 1 minute, we reset the count and timestamp for that IP.
func RateLimitMiddleware(options RateLimitOptions) ghast.Middleware {
	return func(next ghast.Handler) ghast.Handler {
		return ghast.HandlerFunc(func(rw ghast.ResponseWriter, r *ghast.Request) {
			clientIP := r.ClientIP
			entry, exists := rateLimitCollection[clientIP]
			currentTime := time.Now().Unix()

			if exists {
				if currentTime-entry.Timestamp < 60 {
					if entry.Count >= options.RequestsPerMinute {
						rw.Status(429)
						rw.Send([]byte("Too Many Requests"))
						return
					}
					entry.Count++
				} else {
					entry.Count = 1
					entry.Timestamp = currentTime
				}
				rateLimitCollection[clientIP] = entry
			} else {
				rateLimitCollection[clientIP] = rateLimitEntry{
					Count:     1,
					Timestamp: currentTime,
				}
			}
			next.ServeHTTP(rw, r)
		})
	}
}
