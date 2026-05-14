package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"pompidou41/go-shortener/internal/service"

	"pompidou41/go-shortener/internal/utils/validator"
	"regexp"
	"strings"
)

type Handler struct {
	ctx     context.Context
	service service.Service
}

func NewHandler(ctx context.Context, service service.Service) *Handler {
	return &Handler{ctx: ctx, service: service}
}

var shortCodeRegexp = regexp.MustCompile(`^[a-zA-Z0-9]{8}$`)

func (h *Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
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

	code, err := h.service.Shorten(h.ctx, *longUrl)

	if err != nil {
		http.Error(w, "Item not created", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(code))
}

func (h *Handler) LengthenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed, expected: GET, but got: "+r.Method, http.StatusMethodNotAllowed)
		return
	}
	code := strings.TrimPrefix(r.URL.Path, "/")

	if !shortCodeRegexp.MatchString(code) {
		http.NotFound(w, r)
		return
	}

	longUrl, err := h.service.Resolve(h.ctx, code)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longUrl, http.StatusTemporaryRedirect)
}
