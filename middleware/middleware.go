package middleware

import (
	"github.com/gorilla/mux"
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
	span := simpletrace.NewSpan(mux.CurrentRoute(request).GetName())
	for key, value := range h.defaultTags {
		span.Tag(key, value)
	}
	log.Println(span.ParentId)
	span.Start()
	h.next.ServeHTTP(writer, request)
	span.Finalize().Submit(h.client)
}

type Option func(*handler)

func Tags(tags map[string]string) Option {
	return func(h *handler) {
		h.defaultTags = tags
	}
}

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
