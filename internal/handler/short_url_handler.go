package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
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

type ShortURLJsonRequest struct {
	URL string `json:"url"`
}

type ShortURLJsonResponse struct {
	Result string `json:"result"`
}

func (h ShortURLHandler) CreateShortURL(res http.ResponseWriter, req *http.Request) {
	h.HandleShortURL(res, req, false)
}

func (h ShortURLHandler) CreateJSONShortURL(res http.ResponseWriter, req *http.Request) {
	h.HandleShortURL(res, req, true)
}

func (h ShortURLHandler) HandleShortURL(res http.ResponseWriter, req *http.Request, useJSON bool) {
	body, err := io.ReadAll(req.Body)
	req.Body.Close()

	if err != nil {
		http.Error(res, "", http.StatusBadRequest)
		return
	}

	urlFromBody := ""
	if useJSON {
		var jsonRequest ShortURLJsonRequest

		if err := json.Unmarshal(body, &jsonRequest); err != nil {
			http.Error(res, "", http.StatusBadRequest)
			return
		}

		urlFromBody = jsonRequest.URL
	} else {
		urlFromBody = string(body)
	}

	shortURL := ""

	if urlFromBody != "" {
		bytes := make([]byte, 6)
		addURLSuccess := false

		for i := 0; i < 10; i++ {
			rand.Read(bytes)
			shortURL = hex.EncodeToString(bytes)

			if _, ok := h.Urls[shortURL]; !ok {
				(h.Urls)[shortURL] = urlFromBody
				addURLSuccess = true
				break
			}
		}

		if !addURLSuccess {
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Printf("Error: collision was not resolved (url: %s)", urlFromBody)
		}

		baseURL := GetBaseURL(fmt.Sprintf("http://%s", req.Host), h.ResultBaseURL)
		resultURL, ok := url.JoinPath(baseURL, shortURL)

		if ok != nil {
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(ok)
		}

		if useJSON {
			response := ShortURLJsonResponse{
				Result: resultURL,
			}

			jsonResultData, err := json.Marshal(response)

			if err != nil {
				http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Print(err)
			}

			res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(http.StatusCreated)
			res.Write(jsonResultData)
		} else {
			res.WriteHeader(http.StatusCreated)
			fmt.Fprint(res, resultURL)
		}

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
