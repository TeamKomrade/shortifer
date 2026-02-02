package httpserver

import (
	"net/http"

	"github.com/TeamKomrade/shortifer/internal/handler"
)

func CreateServer() {
	var urls = make(map[string]string)

	var urlHandler = handler.ShortUrlHandler{
		Urls: &urls,
	}

	var mux = http.NewServeMux()
	mux.Handle("/", urlHandler)

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		panic(err)
	}
}
