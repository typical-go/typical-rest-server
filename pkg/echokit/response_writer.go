package echokit

import (
	"net/http"
)

type (
	// ResponseWriter http response writer
	ResponseWriter struct {
		Bytes      []byte
		StatusCode int
		header     http.Header
	}
)

var _ http.ResponseWriter = (*ResponseWriter)(nil)

// NewResponseWriter return new instance of ResponseWriter
func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		header: make(http.Header),
	}
}

// Header return http header
func (w *ResponseWriter) Header() http.Header {
	return w.header
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.Bytes = b
	return len(b), nil
}

// WriteHeader sends an HTTP response header with the provided status code
func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

// CopyTo copy to another response writer
func (w *ResponseWriter) CopyTo(rw http.ResponseWriter) {
	rw.Write(w.Bytes)
	rw.WriteHeader(w.StatusCode)
	for k := range w.Header() {
		rw.Header().Add(k, w.Header().Get(k))
	}
}
