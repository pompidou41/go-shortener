package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"pompidou17/shortener/internal/config"
	"pompidou17/shortener/internal/hash"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type Store struct {
	mu      sync.RWMutex
	data    map[string]string
	counter int64
}

var db = Store{
	data: make(map[string]string),
}

var shorCodeRegexp = regexp.MustCompile(`^[a-zA-Z0-9]{8}$`)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func ShortenHandler(conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		salt := conf.SecretSalt

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed, expected: POST, but got: "+r.Method, http.StatusMethodNotAllowed)
			return
		}

		type RequestBody struct {
			URL *string `json:"url"`
		}

		var reqBody RequestBody

		err := json.NewDecoder(r.Body).Decode(&reqBody)

		if err != nil {
			http.Error(w, "Body read err:"+err.Error(), http.StatusBadRequest)
			return
		}

		longUrl := reqBody.URL
		if longUrl == nil || *longUrl == "" {
			http.Error(w, "URL not provided in body or it is empty", http.StatusBadRequest)
			return
		}

		parsed, err := url.ParseRequestURI(*longUrl)

		if err != nil {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		db.mu.Lock()
		defer db.mu.Unlock()

		db.counter++
		counter := db.counter
		shortUrl := hash.EncodeId(counter, salt)

		db.data[shortUrl] = parsed.String()

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortUrl))
	}
}

func LengthenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed, expected: GET, but got: "+r.Method, http.StatusMethodNotAllowed)
		return
	}
	code := strings.TrimPrefix(r.URL.Path, "/")

	if !shorCodeRegexp.MatchString(code) {
		http.NotFound(w, r)
		return
	}

	db.mu.RLock()
	long, exists := db.data[code]
	db.mu.RUnlock()

	if exists {
		http.Redirect(w, r, long, http.StatusTemporaryRedirect)
		return
	}

	http.NotFound(w, r)
}

func main() {
	conf := config.New()

	port := ":" + conf.Port

	mux := http.NewServeMux()

	mux.HandleFunc("/shorten", ShortenHandler(conf))
	mux.HandleFunc("/", LengthenHandler)

	srv := http.Server{
		Addr:           port,
		Handler:        mux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		IdleTimeout:    5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(srv.ListenAndServe())

}
