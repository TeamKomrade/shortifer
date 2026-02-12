package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type ShortURLHandler struct {
	ResultBaseURL string
	Urls          *map[string]string
}

func (h ShortURLHandler) CreateShortURL(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	req.Body.Close()

	if err != nil {
		http.Error(res, "", http.StatusBadRequest)
		return
	}

	urlFromBody := string(body)
	shortURL := ""

	if urlFromBody != "" {
		hash := sha256.Sum256([]byte(urlFromBody))
		shortURL = hex.EncodeToString(hash[:6])

		(*h.Urls)[shortURL] = urlFromBody
		res.WriteHeader(http.StatusCreated)

		baseURL := GetBaseURL(fmt.Sprintf("http://%s", req.Host), h.ResultBaseURL)
		resultURL, ok := url.JoinPath(baseURL, shortURL)

		if ok != nil {
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		fmt.Fprint(res, resultURL)
	} else {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func (h ShortURLHandler) GetFromShortURL(res http.ResponseWriter, req *http.Request) {
	shortURL := path.Base(req.URL.Path)
	longURL, exists := (*h.Urls)[shortURL]

	if !exists {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", longURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func GetBaseURL(defaultURL string, overrideURL string) string {
	if overrideURL != "" {
		return overrideURL
	}

	return defaultURL
}
