package middleware

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	Status int
	Size   int
}

func (r *ResponseWriter) WriteHeader(status int) {
	if r.Status == 0 {
		r.Status = status
		r.ResponseWriter.WriteHeader(status)
	}
}

func (r *ResponseWriter) Write(body []byte) (int, error) {
	if r.Status == 0 {
		r.WriteHeader(http.StatusOK)
	}

	var err error
	r.Size, err = r.ResponseWriter.Write(body)

	return r.Size, err
}
