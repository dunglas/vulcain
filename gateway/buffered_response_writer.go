package gateway

import (
	"bytes"
	"net/http"
)

type bufferedResponseWriter struct {
	originalResponseWriter http.ResponseWriter
	buffer                 *bytes.Buffer
}

func (rw *bufferedResponseWriter) Header() http.Header {
	return rw.originalResponseWriter.Header()
}

func (rw *bufferedResponseWriter) Write(b []byte) (int, error) {
	return rw.buffer.Write(b)
}

func (rw *bufferedResponseWriter) WriteHeader(statusCode int) {
	rw.originalResponseWriter.WriteHeader(statusCode)
}
