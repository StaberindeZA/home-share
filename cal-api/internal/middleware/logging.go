// Package middleware contains logic for middleware
package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rwd := &ResponseWriterDelegator{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rwd, r)

		log.Printf("%s %s completed in %v with %d", r.Method, r.URL.Path, time.Since(start), rwd.statusCode)
	})
}
