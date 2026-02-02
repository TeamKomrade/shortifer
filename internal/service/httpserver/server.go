package httpserver

import (
	"net/http"

	"github.com/TeamKomrade/shortifer/internal/handler"

	chi "github.com/go-chi/chi/v5"
)

func CreateServer() {

	var urls = make(map[string]string)

	var urlHandler = handler.ShortURLHandler{
		Urls: &urls,
	}

	router := chi.NewRouter()

	router.Get("/{shortURL}", urlHandler.GetFromShortURL)
	router.Post("/", urlHandler.CreateShortURL)

	err := http.ListenAndServe(":8080", router)

	if err != nil {
		panic(err)
	}
}
