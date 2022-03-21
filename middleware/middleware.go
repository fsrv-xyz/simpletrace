package middleware

import (
	"fmt"
	"log"
	"net/http"

	"golang.fsrv.services/simpletrace"
)

type handler struct {
	defaultTags map[string]string
	next        http.Handler
	client      *simpletrace.Client
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	span := simpletrace.NewSpan(
		simpletrace.OptionName(request.RequestURI),
		simpletrace.OptionUseKind(simpletrace.KindServer),
		simpletrace.OptionTags(h.defaultTags),
		simpletrace.OptionRemoteEndpoint("client", request.RemoteAddr),
		simpletrace.OptionLocalEndpoint("server", ""),
	)

	// load values from http headers
	parentSpan, err := simpletrace.SpanFromHttpHeader(request)
	if err != nil {
		log.Printf("error while assembling span from headers %+q", err)
	} else {
		span.ParentSpanId = parentSpan.SpanId
		span.TraceId = parentSpan.TraceId
	}

	// set tags for request values
	span.Tag("http.request.host", request.Host)
	span.Tag("http.request.method", request.Method)
	span.Tag("http.request.proto", request.Proto)
	span.Tag("http.request.content.length", request.Header.Get("content-length"))
	span.Tag("http.request.content.type", request.Header.Get("content-type"))
	span.Tag("http.request.user-agent", request.Header.Get("user-agent"))

	log.Printf("trace=%+q", span.TraceId)

	recorder := &ResponseWriter{
		ResponseWriter: writer,
	}

	span.Start()

	// enrich http context with tracing information
	ctx := span.EnrichContext(request.Context(), h.client)

	// forward to handler with enriched context
	h.next.ServeHTTP(recorder, request.WithContext(ctx))

	// set tags for response values
	span.Tag("http.response.code", fmt.Sprintf("%d", recorder.Status))
	span.Tag("http.response.size", fmt.Sprintf("%d", recorder.Size))
	span.Finalize()
	h.client.SubmitAsync(nil, span)
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
