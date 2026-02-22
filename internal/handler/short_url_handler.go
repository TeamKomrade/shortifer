package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
)

type ShortURLHandler struct {
	ResultBaseURL string
	Urls          map[string]string
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
		bytes := make([]byte, 6)
		addUrlSuccess := false

		for i := 0; i < 10; i++ {
			rand.Read(bytes)
			shortURL = hex.EncodeToString(bytes)

			if _, ok := h.Urls[shortURL]; !ok {
				(h.Urls)[shortURL] = urlFromBody
				addUrlSuccess = true
				break
			}
		}

		if !addUrlSuccess {
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Printf("Error: collision was not resolved (url: %s)", urlFromBody)
		}

		res.WriteHeader(http.StatusCreated)

		baseURL := GetBaseURL(fmt.Sprintf("http://%s", req.Host), h.ResultBaseURL)
		resultURL, ok := url.JoinPath(baseURL, shortURL)

		if ok != nil {
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(ok)
		}

		fmt.Fprint(res, resultURL)
	} else {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func (h ShortURLHandler) GetFromShortURL(res http.ResponseWriter, req *http.Request) {
	shortURL := path.Base(req.URL.Path)
	longURL, exists := (h.Urls)[shortURL]

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
