package handler

import (
	"compress/gzip"
	"io"
	"strings"

	http "net/http"
)

type GzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func WithCompression(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.Header.Get("Content-Encoding"), "gzip") {
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
		}

		if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(res, req)
			return
		}

		res.Header().Set("Content-Encoding", "gzip")

		gzwriter := gzip.NewWriter(res)
		defer gzwriter.Close()

		gzw := GzipWriter{
			ResponseWriter: res,
			Writer:         gzwriter,
		}
		h.ServeHTTP(&gzw, req)
	})
}
