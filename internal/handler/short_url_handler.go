package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	pgx "github.com/jackc/pgx/v5"
)

type ShortURLHandler struct {
	ResultBaseURL      string
	DatabaseConnString string
	Urls               map[string]string
	SaveFilePath       string
}

type ShortURLJsonRequest struct {
	URL string `json:"url"`
}

type ShortURLJsonResponse struct {
	Result string `json:"result"`
}

type ShortURLBatchJsonRequestData struct {
	CorrelationId string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ShortURLBatchJsonResponseData struct {
	CorrelationId string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
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

		log.Print("Try save url...")
		h.SaveURLToDB(shortURL, urlFromBody)
		h.SaveURLs()

		baseURL := GetBaseURL(fmt.Sprintf("http://%s", req.Host), h.ResultBaseURL)
		resultURL, ok := url.JoinPath(baseURL, shortURL)

		if ok != nil {
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(ok)
		}

		res.WriteHeader(http.StatusCreated)
		fmt.Fprint(res, resultURL)
	} else {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func (h ShortURLHandler) CreateJSONShortURL(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	req.Body.Close()

	if err != nil {
		http.Error(res, "", http.StatusBadRequest)
		return
	}

	var jsonRequest ShortURLJsonRequest

	if err := json.Unmarshal(body, &jsonRequest); err != nil {
		http.Error(res, "", http.StatusBadRequest)
		return
	}

	urlFromBody := jsonRequest.URL

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

		log.Print("Try save url...")
		h.SaveURLToDB(shortURL, urlFromBody)
		h.SaveURLs()

		baseURL := GetBaseURL(fmt.Sprintf("http://%s", req.Host), h.ResultBaseURL)
		resultURL, ok := url.JoinPath(baseURL, shortURL)

		if ok != nil {
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(ok)
		}

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
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func (h ShortURLHandler) CreateJSONShortURLFromBatch(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	req.Body.Close()

	if err != nil {
		http.Error(res, "", http.StatusBadRequest)
		return
	}

	var jsonRequest []ShortURLBatchJsonRequestData

	if err := json.Unmarshal(body, &jsonRequest); err != nil {
		http.Error(res, "", http.StatusBadRequest)
		return
	}

	resultBatch := make([]ShortURLBatchJsonResponseData, len(jsonRequest))

	for index, value := range jsonRequest {
		var shortURL string
		bytes := make([]byte, 6)
		addURLSuccess := false

		for i := 0; i < 10; i++ {
			rand.Read(bytes)
			shortURL = hex.EncodeToString(bytes)

			if _, ok := h.Urls[shortURL]; !ok {
				(h.Urls)[shortURL] = value.OriginalURL
				addURLSuccess = true
				break
			}
		}

		if !addURLSuccess {
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Printf("Error: collision was not resolved (url: %s)", value.OriginalURL)
		}

		log.Print("Try save url...")
		h.SaveURLToDB(shortURL, value.OriginalURL)
		h.SaveURLs()

		baseURL := GetBaseURL(fmt.Sprintf("http://%s", req.Host), h.ResultBaseURL)
		resultURL, ok := url.JoinPath(baseURL, shortURL)

		if ok != nil {
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(ok)
		}

		resultBatch[index].CorrelationId = value.CorrelationId
		resultBatch[index].ShortURL = resultURL
	}

	jsonResultData, err := json.Marshal(resultBatch)

	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(jsonResultData)
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

func (h ShortURLHandler) SaveURLs() {
	file, err := os.OpenFile(h.SaveFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Print(err)
		return
	}
	defer file.Close()

	jsonData, err := json.Marshal(h.Urls)

	if err != nil {
		log.Print(err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		log.Print(err)
	}
	defer file.Close()
}

func (h ShortURLHandler) SaveURLToDB(shortURL string, originalURL string) {
	conn, err := pgx.Connect(context.Background(), h.DatabaseConnString)
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "INSERT INTO short_url (short_url, original_url) VALUES ($1, $2)", shortURL, originalURL)
	if err != nil {
		log.Print(err)
	}

	log.Print("Saved: original url: ")
}
