package service

import (
	"net/http"
	"time"
)

type Interceptor func(http.Handler) http.Handler
type Interceptors []Interceptor

func (ins Interceptors) Apply(handler http.Handler) http.Handler {
	if len(ins) == 0 {
		return handler
	}
	return ins[1:].Apply(ins[0](handler))
}

func (c *HttpController) Logging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func(start time.Time) {
			requestID := w.Header().Get("X-Request-Id")
			if requestID == "" {
				requestID = "unknown"
			}
			c.Logger.Println(requestID, req.Method, req.URL.Path, req.RemoteAddr, req.UserAgent(), time.Since(start))
		}(time.Now())
		handler.ServeHTTP(w, req)
	})
}

func (c *HttpController) Tracing(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestID := req.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = c.NextRequestID()
		}
		w.Header().Set("X-Request-Id", requestID)
		handler.ServeHTTP(w, req)
	})
}

var (
	_ Interceptor = (&HttpController{}).Logging
	_ Interceptor = (&HttpController{}).Tracing
)
