package httpserver

import (
	"net/http"
	"os"

	"github.com/TeamKomrade/shortifer/internal/handler"

	cfg "github.com/TeamKomrade/shortifer/internal/config"
	chi "github.com/go-chi/chi/v5"
)

func CreateServer(flags cfg.StartupFlags) error {

	hostURL := ":8080"

	urls := make(map[string]string)
	urlHandler := handler.ShortURLHandler{
		Urls:          urls,
		ResultBaseURL: flags.ResultBaseURL,
	}

	envResultAddress := os.Getenv("BASE_URL")
	if envResultAddress != "" {
		urlHandler.ResultBaseURL = envResultAddress
	}

	router := chi.NewRouter()

	router.Get("/{shortURL}", urlHandler.GetFromShortURL)
	router.Post("/", urlHandler.CreateShortURL)

	if flags.BaseURL != "" {
		hostURL = flags.BaseURL
	}

	envHostAddress := os.Getenv("SERVER_ADDRESS")
	if envHostAddress != "" {
		hostURL = envHostAddress
	}

	err := http.ListenAndServe(hostURL, router)

	if err != nil {
		return err
	}

	return nil
}
