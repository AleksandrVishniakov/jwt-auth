package handlers

import (
	"log/slog"
	"net"
	"net/http"
	"time"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	status int
}

func (w *responseWriterWrapper) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func Logger(
	log *slog.Logger,
) func(next http.Handler) http.Handler {
	const src = "handlers.Logger"
	log = log.With(slog.String("src", src))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &responseWriterWrapper{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(wrapped, r)

			attrs := []any{
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", wrapped.status),
				slog.String("duration", time.Since(start).String()),
			}

			if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
				attrs = append(attrs, slog.String("ip", host))
			}

			if wrapped.status >= http.StatusBadRequest {
				log.Error("request failed", attrs...)
			} else {
				log.Info("request completed", attrs...)
			}
		})
	}
}
