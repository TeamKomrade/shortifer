package handler

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"path"
)

type ShortURLHandler struct {
	Urls *map[string]string
}

func (h ShortURLHandler) CreateShortURL(res http.ResponseWriter, req *http.Request) {
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
			shortURL = hex.EncodeToString(hash[:4])

			(*h.Urls)[shortURL] = url
			res.WriteHeader(http.StatusCreated)
			fmt.Fprintf(res, "http://%s/%s", req.Host, shortURL)
		} else {
			http.Error(res, "", http.StatusBadRequest)
		}
	}
}

func (h ShortURLHandler) GetFromShortURL(res http.ResponseWriter, req *http.Request) {
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
}
