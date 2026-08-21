package middleware

import "net/http"

type ResponseWriterDelegator struct {
	http.ResponseWriter
	statusCode int
}

func (rwd *ResponseWriterDelegator) WriteHeader(code int) {
	rwd.statusCode = code
	rwd.ResponseWriter.WriteHeader(code)
}

func (rwd *ResponseWriterDelegator) Unwrap() http.ResponseWriter {
	return rwd.ResponseWriter
}
