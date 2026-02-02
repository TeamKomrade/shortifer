package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func InitShortURLHandler() http.Handler {
	var urls = make(map[string]string)
	var urlHandler = ShortURLHandler{
		Urls: &urls,
	}
	return urlHandler
}

func Test_Shorting(t *testing.T) {
	tests := []struct {
		name               string
		address            string
		expectedMaxLength  int
		expectedStatusCode int
	}{
		{
			name:               "google.com",
			address:            "http://google.com",
			expectedMaxLength:  8,
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "git long commit",
			address:            "https://github.com/TeamKomrade/shortifer/commit/95cf9e311c15c162862adcbb2f710987ee8be04e",
			expectedMaxLength:  8,
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "empty address",
			address:            "",
			expectedMaxLength:  8,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.address))
			rec := httptest.NewRecorder()

			ShortURLHandler := InitShortURLHandler()
			ShortURLHandler.ServeHTTP(rec, request)

			result := rec.Result()

			defer result.Body.Close()
			body, err := io.ReadAll(result.Body)

			require.NoError(t, err)

			shortURL := path.Base(string(body))

			assert.Equal(t, test.expectedStatusCode, result.StatusCode)

			if test.expectedStatusCode != http.StatusBadRequest {
				assert.NotEmpty(t, body)
				assert.LessOrEqual(t, len(shortURL), test.expectedMaxLength)
			}
		})
	}
}
