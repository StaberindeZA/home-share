// Package middleware contains logic for middleware
package middleware

import (
	"log/slog"
	"net/http"
	"slices"
	"time"
)

func Logger(next http.Handler) http.Handler {
	debugOnlyPaths := []string{"/metrics", "/livez", "/readyz"}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rwd := &ResponseWriterDelegator{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rwd, r)

		if slices.Contains(debugOnlyPaths, r.URL.Path) {
			slog.Debug("request completed", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start), "statusCode", rwd.statusCode)
		} else {
			slog.Info("request completed", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start), "statusCode", rwd.statusCode)
		}
	})
}
