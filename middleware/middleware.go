package middleware

import (
	"fmt"
	"golang.fsrv.services/simpletrace"
	"log"
	"net/http"
)

type handler struct {
	defaultTags map[string]string
	next        http.Handler
	client      *simpletrace.Client
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	span := simpletrace.NewSpan(
		request.RequestURI,
		simpletrace.UseKind(simpletrace.KindServer),
		simpletrace.Tags(h.defaultTags),
		simpletrace.RemoteEndpoint("client", request.RemoteAddr),
		simpletrace.LocalEndpoint("server", ""),
	)
	span.Tag("http.request.host", request.Host)
	span.Tag("http.request.method", request.Method)
	span.Tag("http.request.proto", request.Proto)
	span.Tag("http.request.user-agent", request.Header.Get("user-agent"))

	log.Printf("trace=%+q", span.TraceId)

	recorder := &ResponseWriter{
		ResponseWriter: writer,
	}

	span.Start()
	h.next.ServeHTTP(recorder, request)
	span.Tag("http.response.code", fmt.Sprintf("%d", recorder.Status))
	span.Tag("http.response.size", fmt.Sprintf("%d", recorder.Size))
	span.Finalize().Submit(h.client)
}

type Option func(*handler)

// Tags - assign default tags to the middleware
func Tags(tags map[string]string) Option {
	return func(h *handler) {
		h.defaultTags = tags
	}
}

// NewMiddleware - create a new middleware with options
func NewMiddleware(c *simpletrace.Client, options ...Option) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := &handler{
			next:   next,
			client: c,
		}
		for _, option := range options {
			option(h)
		}
		return h
	}
}
