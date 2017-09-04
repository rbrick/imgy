package main

import (
	"log"
	"net/http"

	"github.com/fatih/color"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (l *loggingResponseWriter) WriteHeader(status int) {
	l.status = status
	l.ResponseWriter.WriteHeader(status)
}

type logHandler struct {
	h http.Handler
}

func (l *logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lrw := &loggingResponseWriter{w, http.StatusOK}
	l.h.ServeHTTP(lrw, r)

	status := color.New(colorForStatus(lrw.status)).Sprint(lrw.status)
	method := color.New(colorForMethod(r.Method)).Sprint(r.Method)

	// [Imgy] (GET | 200) - /
	log.Printf("[Imgy] (%s | %s) - %s | %s", method, status, r.RequestURI, r.RemoteAddr)
}

func NewLogHandler(h http.Handler) http.Handler {
	return &logHandler{h}
}

func colorForStatus(code int) color.Attribute {
	switch {
	case code >= 200 && code < 300:
		return color.BgGreen
	case code >= 300 && code < 400:
		return color.BgCyan
	case code >= 400 && code < 500:
		return color.BgYellow
	default:
		return color.BgRed
	}
}

func colorForMethod(method string) color.Attribute {
	switch method {
	case "GET":
		return color.BgGreen
	case "POST":
		return color.BgCyan
	case "PUT":
		return color.BgYellow
	case "DELETE":
		return color.BgRed
	case "PATCH":
		return color.BgBlue
	case "HEAD":
		return color.BgMagenta
	case "OPTIONS":
		return color.BgWhite
	default:
		return color.Reset
	}
}
