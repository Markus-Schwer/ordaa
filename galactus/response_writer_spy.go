package main

import (
	"net/http"
)

type ResponseWriterSpy struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriterSpy(w http.ResponseWriter) *ResponseWriterSpy {
	return &ResponseWriterSpy{ResponseWriter: w, statusCode: 200}
}

func (w *ResponseWriterSpy) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *ResponseWriterSpy) Write(body []byte) (int, error) {
	return w.ResponseWriter.Write(body)
}
