package gateway

// Adapted from https://github.com/gabesullice/hades/blob/master/lib/server/response_buffer.go
// Copyright (c) 2019 Gabriel Sullice
// MIT License

import (
	"bytes"
	"io"
	"net/http"
)

type bufferedResponseWriter struct {
	originalResponseWriter http.ResponseWriter
	buffer                 *bytes.Buffer
	statusCode             int
	body                   []byte
}

func newBufferedResponseWriter(rw http.ResponseWriter) *bufferedResponseWriter {
	var buffer bytes.Buffer
	return &bufferedResponseWriter{originalResponseWriter: rw, buffer: &buffer}
}

func (rw *bufferedResponseWriter) Header() http.Header {
	return rw.originalResponseWriter.Header()
}

func (rw *bufferedResponseWriter) Write(b []byte) (int, error) {
	return rw.buffer.Write(b)
}

// WriteHeader only stores the status code. Headers will be really written only send() will be called.
// This trick allow the reverse proxy to add Link preload headers.
func (rw *bufferedResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
}

func (rw *bufferedResponseWriter) bodyContent() []byte {
	if rw.body == nil {
		rw.body = rw.buffer.Bytes()
	}

	return rw.body
}

func (rw *bufferedResponseWriter) send() {
	rw.originalResponseWriter.WriteHeader(rw.statusCode)
	io.Copy(rw.originalResponseWriter, bytes.NewReader(rw.bodyContent()))
	rw.buffer.Reset()
	rw.body = nil
}
