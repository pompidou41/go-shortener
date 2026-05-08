package handler

import (
	"encoding/json"
	"net/http"
	"pompidou41/go-shortener/internal/config"
	"pompidou41/go-shortener/internal/hash"
	"pompidou41/go-shortener/internal/storage"
	"pompidou41/go-shortener/internal/utils/validator"
	"regexp"
	"strings"
)

type Handler struct {
	store *storage.Store
	conf  *config.Config
}

func NewHandler(store *storage.Store, conf *config.Config) *Handler {
	return &Handler{store: store, conf: conf}
}

var shorCodeRegexp = regexp.MustCompile(`^[a-zA-Z0-9]{8}$`)

func (h *Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	salt := h.conf.SecretSalt

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

	parsingErr := validator.ValidateURL(*longUrl)

	if parsingErr != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	h.store.Mu.Lock()
	defer h.store.Mu.Unlock()

	h.store.Counter++
	counter := h.store.Counter
	shortUrl := hash.EncodeId(counter, salt)

	h.store.Data[shortUrl] = *longUrl

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortUrl))
}

func (h *Handler) LengthenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed, expected: GET, but got: "+r.Method, http.StatusMethodNotAllowed)
		return
	}
	code := strings.TrimPrefix(r.URL.Path, "/")

	if !shorCodeRegexp.MatchString(code) {
		http.NotFound(w, r)
		return
	}

	h.store.Mu.RLock()
	long, exists := h.store.Data[code]
	h.store.Mu.RUnlock()

	if exists {
		http.Redirect(w, r, long, http.StatusTemporaryRedirect)
		return
	}

	http.NotFound(w, r)
}
