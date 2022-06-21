package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fsrv-xyz/simpletrace"
)

type handler struct {
	defaultTags       map[string]string
	next              http.Handler
	spanSubmitChannel chan *simpletrace.Span
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if h.spanSubmitChannel == nil {
		panic("spanSubmitChannel may not be empty")
	}
	span := simpletrace.NewSpan(
		simpletrace.OptionName(request.RequestURI),
		simpletrace.OptionOfKind(simpletrace.KindServer),
		simpletrace.OptionTags(h.defaultTags),
		simpletrace.OptionRemoteEndpoint("client", request.RemoteAddr),
		simpletrace.OptionLocalEndpoint("server", ""),
	)

	// load values from http headers
	parentSpan, err := simpletrace.SpanFromHttpHeader(request)
	if err != nil {
		log.Printf("cannot assemble span from headers %+q", err)
	} else {
		span.Use(simpletrace.OptionFromParent(parentSpan.SpanId))
		span.Use(simpletrace.OptionTraceID(parentSpan.TraceId))
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

	// forward to handler with enriched context
	h.next.ServeHTTP(recorder, request.WithContext(span.EnrichContext(request.Context())))

	// set tags for response values
	span.Tag("http.response.code", fmt.Sprintf("%d", recorder.Status))
	span.Tag("http.response.size", fmt.Sprintf("%d", recorder.Size))
	span.Finalize()
	h.spanSubmitChannel <- span
}

type Option func(*handler)

// Tags - assign default tags to the middleware
func Tags(tags map[string]string) Option {
	return func(h *handler) {
		h.defaultTags = tags
	}
}

// NewMiddleware - create a new middleware with options
func NewMiddleware(c chan *simpletrace.Span, options ...Option) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := &handler{
			next:              next,
			spanSubmitChannel: c,
		}
		for _, option := range options {
			option(h)
		}
		return h
	}
}
