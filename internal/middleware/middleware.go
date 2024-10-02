package middleware

import (
	"log"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			x := middlewares[i]
			next = x(next)
		}
		return next
	}
}

type WrappedWritter struct {
	http.ResponseWriter
	statusCode int
}

func (writer *WrappedWritter) WriteHeader(statusCode int) {
	writer.ResponseWriter.WriteHeader(statusCode)
	writer.statusCode = statusCode
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		wrapped := &WrappedWritter{
			ResponseWriter: writer,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapped, request)
		log.Println("INFO", wrapped.statusCode, request.Method, request.URL.Path, time.Since(start))
	})
}
