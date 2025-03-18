package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	//"github.com/lithammer/shortuuid/v4"
	"jaxon.app/jaxon/internal/auth/sessions"
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

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &WrappedWritter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		// reqId := shortuuid.New()
		next.ServeHTTP(wrapped, r)
		slog.Info(r.URL.Path,
			"method", r.Method,
			"status", wrapped.statusCode,
			"time", time.Since(start))
	})
}

func EnsureAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// before
		cookie, err := r.Cookie("trbl_session")
		if err != nil {
			// user is not authenticated, need to redirect
			// have to handle regular redirect and hx redirect
			slog.Warn("No session cookie found")
			_, ok := r.Header["Hx-Request"]
			if ok {
				w.Header().Add("Hx-Redirect", "/auth/login")
				w.WriteHeader(http.StatusOK)
			} else {
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			}
			return
		}

		db := r.Context().Value("db").(*sqlx.DB)
		session, err := sessions.Fetch(cookie.Value, db)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if session == nil {
			slog.Warn("No session found for given cookie. Suspicious", "cookie", cookie.Value)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if session.IsExpired() {
			slog.Info("Session is expired", "session", session)
			_, ok := r.Header["Hx-Request"]
			if ok {
				w.Header().Add("Hx-Redirect", "/auth/login")
				w.WriteHeader(http.StatusOK)
			} else {
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			}
			return
		}

		ctx := context.WithValue(r.Context(), "userId", session.UserId)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
