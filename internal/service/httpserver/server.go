package httpserver

import (
	"net/http"

	"github.com/TeamKomrade/shortifer/internal/handler"

	cfg "github.com/TeamKomrade/shortifer/internal/config"
	chi "github.com/go-chi/chi/v5"
)

func CreateServer(flags cfg.StartupFlags) {
	hostURL := ":8080"

	var urls = make(map[string]string)
	var urlHandler = handler.ShortURLHandler{
		Urls:          &urls,
		ResultBaseURL: flags.ResultBaseURL,
	}

	router := chi.NewRouter()

	router.Get("/{shortURL}", urlHandler.GetFromShortURL)
	router.Post("/", urlHandler.CreateShortURL)

	if flags.BaseURL != "" {
		hostURL = flags.BaseURL
	}

	err := http.ListenAndServe(hostURL, router)

	if err != nil {
		panic(err)
	}
}
