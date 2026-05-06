package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"pompidou17/shortener/internal/config"
	"pompidou17/shortener/internal/hash"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

type url struct {
	id       int64
	shortUrl string
	longUrl  string
}

var db []url

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Method not allowed, expected: POST, but got:", r.Method)
		http.Error(w, "Method not allowed, expected: POST, but got: "+r.Method, http.StatusMethodNotAllowed)
	} else {
		type RequestBody struct {
			URL *string `json:"url"`
		}
		// TODO: refactor secrets
		salt := config.New().SecretSalt

		jsonBlob, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Body read err:", err)
			http.Error(w, "Body read err:"+err.Error(), http.StatusBadRequest)
			return
		}

		var reqBody RequestBody
		err = json.Unmarshal(jsonBlob, &reqBody)

		if err != nil {
			log.Println("Body parse err:", err)
			http.Error(w, "Body parse err:"+err.Error(), http.StatusBadRequest)
			return
		}

		longUrl := reqBody.URL
		if longUrl == nil || *reqBody.URL == "" {
			log.Println("URL not provided in body or it is empty")
			http.Error(w, "URL not provided in body or it is empty", http.StatusBadRequest)
			return
		} else {
			id := int64(len(db) + 1)
			shortUrl := hash.EncodeId(id, salt)

			url := url{
				id:       id,
				longUrl:  *longUrl,
				shortUrl: shortUrl,
			}

			db = append(db, url)
			log.Println(longUrl, url, db)

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(shortUrl))
			return
		}
	}
}

func decodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("Method not allowed, expected: GET, but got:", r.Method)
		http.Error(w, "Method not allowed, expected: GET, but got: "+r.Method, http.StatusMethodNotAllowed)
		return
	}
	code := strings.TrimPrefix(r.URL.Path, "/")

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]{8}$`, code)
	if !matched {
		http.NotFound(w, r)
		return
	}

	// TODO: optimize to O(1)
	for _, url := range db {
		if url.shortUrl == code {
			log.Printf("Redirecting %s -> %s", code, url.longUrl)
			http.Redirect(w, r, url.longUrl, http.StatusMovedPermanently)
			return
		}
	}
	log.Printf("Code %s not found", code)
	http.NotFound(w, r)
}

func main() {
	conf := config.New()

	port := ":" + conf.Port

	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", decodeHandler)
	log.Fatal((http.ListenAndServe(port, nil)))
}
