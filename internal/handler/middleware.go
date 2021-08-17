package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog"
)

type LogFormatter struct {
	*zerolog.Logger
}

func (l *LogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	logger := l.With().
		Str("req", fmt.Sprintf("%s://%s%s %s", scheme, r.Host, r.RequestURI, r.Proto)).
		Str("from", r.RemoteAddr).
		Logger()
	return &LogEntry{&logger}
}

type LogEntry struct {
	*zerolog.Logger
}

func (l *LogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Info().
		Int("status", status).
		Int("bytes", bytes).
		Str("elapsed", elapsed.String()).
		Msg("Request handled")
}

func (l *LogEntry) Panic(v interface{}, stack []byte) {
	l.Info().
		Interface("panic", v).
		Bytes("stack", stack).
		Msg("Panic handled")
}

func JWT(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/auth":
			default:
				authHeader := r.Header.Get("Authorization")
				if len(authHeader) == 0 {
					writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
					return
				}
				h := strings.SplitN(authHeader, " ", 2)
				if len(h) != 2 {
					writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
					return
				}
				if strings.ToLower(h[0]) != "bearer" {
					writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
					return
				}
				_, err := jwt.ParseString(h[1], jwt.WithVerify(jwa.HS256, secret), jwt.WithValidate(true))
				if err != nil {
					writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
