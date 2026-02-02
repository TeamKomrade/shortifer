package handler

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"path"
)

type ShortURLHandler struct {
	Urls *map[string]string
}

func (h ShortURLHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		shortURL := path.Base(req.URL.Path)
		longURL, exists := (*h.Urls)[shortURL]

		if !exists {
			http.Error(res, "URL not found", http.StatusBadRequest)
			return
		}

		res.Header().Set("Location", longURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}

	if req.Method == http.MethodPost {
		body, err := io.ReadAll(req.Body)
		req.Body.Close()

		if err != nil {
			http.Error(res, "", http.StatusBadRequest)
			return
		}

		url := string(body)
		shortURL := ""

		if url != "" {
			hash := md5.Sum([]byte(url))
			shortURL = string(hash[:8])

			(*h.Urls)[shortURL] = url
			res.WriteHeader(http.StatusCreated)
			fmt.Fprintf(res, "http://%s/%s", req.Host, shortURL)
		} else {
			http.Error(res, "", http.StatusBadRequest)
		}
	}
}
