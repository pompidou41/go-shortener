package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"pompidou17/shortener/internal/config"
	"pompidou17/shortener/internal/hash"
	"strings"
	"testing"
)

func TestMainHandlers(t *testing.T) {
	t.Run("POST /shorten", func(t *testing.T) {
		conf := config.Config{
			SecretSalt: "salt",
		}

		cases := []struct {
			name           string
			method         string
			body           any
			expectedStatus int
			expectedBody   string
		}{
			{
				name:           "Success",
				method:         http.MethodPost,
				body:           map[string]string{"url": "https://example.com/"},
				expectedStatus: http.StatusCreated,
				expectedBody:   hash.EncodeId(1, "salt"),
			},
			{
				name:           "Invalid method",
				method:         http.MethodGet,
				body:           nil,
				expectedStatus: http.StatusMethodNotAllowed,
				expectedBody:   "Method not allowed",
			},
			{
				name:           "Empty URL in body",
				method:         http.MethodPost,
				body:           map[string]string{"url": ""},
				expectedStatus: http.StatusBadRequest,
				expectedBody:   "URL not provided",
			},
			{
				name:           "Invalid URL",
				method:         http.MethodPost,
				body:           map[string]string{"url": "invalid-url"},
				expectedStatus: http.StatusBadRequest,
				expectedBody:   "Invalid URL",
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				db = Store{
					data: make(map[string]string),
				}

				var buf bytes.Buffer

				if tc.body != nil {
					err := json.NewEncoder(&buf).Encode(tc.body)
					if err != nil {
						t.Fatalf("Encode body error: %v", err)
					}
				}

				req := httptest.NewRequest(tc.method, "/shorten", &buf)
				w := httptest.NewRecorder()

				handler := ShortenHandler(&conf)
				handler(w, req)

				res := w.Result()
				defer res.Body.Close()

				if res.StatusCode != tc.expectedStatus {
					t.Fatalf("Expected status %d, got %d", tc.expectedStatus, res.StatusCode)
				}

				data, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("Read body error: %v", err)
				}

				if !strings.Contains(string(data), tc.expectedBody) {
					t.Errorf("Expected body to contain %q, got %q", tc.expectedBody, string(data))
				}
			})
		}
	})

	t.Run("GET /", func(t *testing.T) {
		db = Store{
			data: map[string]string{
				"dsF521cZ": "https://example.com/",
				"g638sAGc": "https://onemore.site/page/2",
			},
			counter: 2,
		}

		cases := []struct {
			name           string
			method         string
			param          string
			expectedStatus int
			expectedHeader map[string]string
		}{
			{
				name:           "Success redirect",
				method:         http.MethodGet,
				param:          "g638sAGc",
				expectedStatus: http.StatusTemporaryRedirect,
				expectedHeader: map[string]string{"Location": db.data["g638sAGc"]},
			},
			{
				name:           "Invalid method",
				method:         http.MethodPost,
				param:          "g638sAGc",
				expectedStatus: http.StatusMethodNotAllowed,
			},
			{
				name:           "Invalid code",
				method:         http.MethodGet,
				param:          "123abc",
				expectedStatus: http.StatusNotFound,
			},
			{
				name:           "Empty code",
				method:         http.MethodGet,
				param:          "",
				expectedStatus: http.StatusNotFound,
			},
			{
				name:           "Code with invalid symbols",
				method:         http.MethodGet,
				param:          "image.png",
				expectedStatus: http.StatusNotFound,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest(tc.method, "/"+tc.param, nil)
				w := httptest.NewRecorder()

				LengthenHandler(w, req)

				res := w.Result()

				if res.StatusCode != tc.expectedStatus {
					t.Fatalf("Expected status %d, got %d", tc.expectedStatus, res.StatusCode)
				}

				for key, expectedValue := range tc.expectedHeader {
					actualValue := res.Header.Get(key)

					if actualValue != expectedValue {
						t.Errorf("Expected header %s to be %q, got %q", key, expectedValue, actualValue)
					}
				}
			})
		}
	})
}
