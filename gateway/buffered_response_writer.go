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

func (rw *bufferedResponseWriter) WriteHeader(statusCode int) {
	rw.originalResponseWriter.WriteHeader(statusCode)
}

func (rw *bufferedResponseWriter) bodyContent() []byte {
	if rw.body == nil {
		rw.body = rw.buffer.Bytes()
	}

	return rw.body
}

func (rw *bufferedResponseWriter) send() {
	io.Copy(rw.originalResponseWriter, bytes.NewReader(rw.bodyContent()))
	rw.buffer.Reset()
	rw.body = nil
}
