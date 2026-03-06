package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/swagger") {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		traceID := chiMiddleware.GetReqID(r.Context())

		var reqBody []byte
		if r.Body != nil {
			reqBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		// Wrap Response Writer
		rw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
			body:           bytes.NewBuffer(nil),
		}

		next.ServeHTTP(rw, r)

		latency := time.Since(start)

		// Structured Logging
		slog.LogAttrs(r.Context(), slog.LevelInfo, "HTTP Request",
			slog.String("trace_id", traceID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("ip", r.RemoteAddr),
			slog.Int("status", rw.status),
			slog.Duration("latency", latency),
			slog.String("request_body", string(truncate(reqBody, 1024))),
			slog.String("response_body", string(truncate(rw.body.Bytes(), 1024))),
		)
	})
}

func truncate(b []byte, n int) []byte {
	if len(b) > n {
		return b[:n]
	}
	return b
}
