package flusher_response_writer

import (
	"context"
	"log"
	"net/http"
	"reflect"
)

type Config struct{}

func CreateConfig() *Config {
	return &Config{}
}

type Plugin struct {
	next http.Handler
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &Plugin{
		next: next,
	}, nil
}

type responseWriter struct {
	http.ResponseWriter
}

func (w *responseWriter) Flush() {
	log.Println("Flushing...")
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (plugin *Plugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var respCodeOverrideWriter = &responseWriter{ResponseWriter: rw}

	// Some test code here just to showcase the problem:
	var flusher http.Flusher = respCodeOverrideWriter // This compiles and works fine. responseWriter implements http.Flusher
	log.Println("responseWriter is http.Flusher. We can call Flush manually:")
	flusher.Flush()

	var httpResponseWriter http.ResponseWriter = respCodeOverrideWriter // This also compiles and works fine
	_, ok := httpResponseWriter.(http.Flusher)
	log.Println("Is responseWriter http.Flusher?", ok) // false?? Type conversion failed. So responseWriter is not a http.Flusher after all?

	log.Println("Type of responseWriter: ", reflect.TypeOf(httpResponseWriter)) // stdlib._net_http_ResponseWriter. Somehow the underlying type got changed?

	// This never calls my Flush implementation:
	plugin.next.ServeHTTP(respCodeOverrideWriter, req)
}
