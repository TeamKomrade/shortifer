package httpserver

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/TeamKomrade/shortifer/internal/handler"

	cfg "github.com/TeamKomrade/shortifer/internal/config"
	chi "github.com/go-chi/chi/v5"
)

func CreateServer(flags cfg.StartupFlags) error {

	hostURL := ":8080"
	router := chi.NewRouter()

	router.Use(handler.WithLog)
	router.Use(handler.WithCompression)

	urlHandler := handler.ShortURLHandler{}

	jsonFilePath := os.Getenv("FILE_STORAGE_PATH")
	if jsonFilePath == "" {
		if flags.JsonFilePath != "" {
			jsonFilePath = flags.JsonFilePath
		} else {
			jsonFilePath = "C:\\shorter\\urls.json"
		}
	}
	SetupShortURLHandler(urlHandler, router, flags, jsonFilePath)

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

func SetupShortURLHandler(shortURLHandler handler.ShortURLHandler, router chi.Router, flags cfg.StartupFlags, jsonFilePath string) {

	urls := make(map[string]string)
	if jsonFilePath != "" {
		data, err := os.ReadFile(jsonFilePath)

		if err == nil {
			if err := json.Unmarshal(data, &urls); err != nil {
				log.Print(err)
			}
		} else {
			log.Print(err)
		}
	}

	urlHandler := handler.ShortURLHandler{
		Urls:          urls,
		ResultBaseURL: flags.ResultBaseURL,
		SaveFilePath:  jsonFilePath,
	}

	envResultAddress := os.Getenv("BASE_URL")
	if envResultAddress != "" {
		urlHandler.ResultBaseURL = envResultAddress
	}

	router.Get("/{shortURL}", urlHandler.GetFromShortURL)
	router.Post("/", urlHandler.CreateShortURL)
	router.Post("/api/shorten", urlHandler.CreateJSONShortURL)
}
