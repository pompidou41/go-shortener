package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"pompidou17/shortener/internal/config"
	"pompidou17/shortener/internal/hash"

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

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(shortUrl))
		return
	}
}

func main() {
	conf := config.New()

	port := ":" + conf.Port

	http.HandleFunc("/shorten", shortenHandler)
	log.Fatal((http.ListenAndServe(port, nil)))
}
