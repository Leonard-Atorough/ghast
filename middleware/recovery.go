package middleware

import (
	"log"

	"github.com/Leonard-Atorough/ghast"
)

// RecoveryMiddleware is a middleware that recovers from panics in handlers and returns a 500 error.
type Options struct {
	Log    bool        // Whether to log the panic error (default: true)
	Logger *log.Logger // Optional custom logger (default: standard logger)
}

// RecoveryMiddleware creates a RecoveryMiddleware with the given options.
func RecoveryMiddleware(opts Options) ghast.Middleware {
	return func(next ghast.Handler) ghast.Handler {
		return ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
			defer func() {
				if err := recover(); err != nil {
					if opts.Log {
						if opts.Logger != nil {
							opts.Logger.Printf("Panic recovered: %v", err)
						} else {
							log.Printf("Panic recovered: %v", err)
						}
					}
					w.JSON(500, map[string]string{"error": "Internal Server Error"})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
