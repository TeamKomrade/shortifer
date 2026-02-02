package handler

import (
	"fmt"
	"io"
	"net/http"
	"path"
)

type ShortUrlHandler struct {
	Urls *map[string]string
}

func (h ShortUrlHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		shortUrl := path.Base(req.URL.Path)
		longUrl, exists := (*h.Urls)[shortUrl]

		if !exists {
			http.Error(res, "URL not found", http.StatusBadRequest)
			return
		}

		res.WriteHeader(http.StatusTemporaryRedirect)
		res.Header().Add("Location", longUrl)

	}

	if req.Method == http.MethodPost {
		body, err := io.ReadAll(req.Body)
		req.Body.Close()

		if err != nil {
			http.Error(res, "", http.StatusBadRequest)
			return
		}

		url := string(body)
		shortUrl := ""

		if url != "" {
			for i, chr := range url {
				if i%2 == 0 {
					if chr != '/' && chr != '\\' && chr != '.' && chr != ':' {
						shortUrl += string(chr)
					}
				}
			}

			(*h.Urls)[shortUrl] = url
			res.WriteHeader(http.StatusCreated)
			fmt.Fprint(res, shortUrl)
		} else {
			http.Error(res, "", http.StatusBadRequest)
		}
	}
}
