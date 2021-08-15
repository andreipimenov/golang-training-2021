package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
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
