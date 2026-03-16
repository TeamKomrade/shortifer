package handler

import (
	"compress/gzip"
	"io"
	"log"
	"strings"

	http "net/http"
)

type GzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func WithCompression(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(res, req)
			return
		}

		gzw, err := gzip.NewWriterLevel(res, gzip.BestSpeed)

		if err != nil {
			log.Printf("Gzip error: %v", err)
			return
		}
		defer gzw.Close()

		gzr, err := gzip.NewReader(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gzr.Close()

		body, err := io.ReadAll(gzr)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Body = io.NopCloser(strings.NewReader(string(body)))

		res.Header().Set("Content-Encoding", "gzip")
		h.ServeHTTP(GzipWriter{ResponseWriter: res, Writer: gzw}, req)
	})
}
