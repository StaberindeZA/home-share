package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rwd := &responseWriterDelegator{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rwd, r)

		log.Printf("%s %s completed in %v with %d", r.Method, r.URL.Path, time.Since(start), rwd.statusCode)
	})
}

// You must include this helper struct for the function to work
type responseWriterDelegator struct {
	http.ResponseWriter
	statusCode int
}

func (rwd *responseWriterDelegator) WriteHeader(code int) {
	rwd.statusCode = code
	rwd.ResponseWriter.WriteHeader(code)
}

func (rwd *responseWriterDelegator) Unwrap() http.ResponseWriter {
	return rwd.ResponseWriter
}
