package handlers

import (
	"encoding/json"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/config"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/models"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/usecase"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func HandlerPost(res http.ResponseWriter, req *http.Request, useCase usecase.ShortURLUseCase) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	if len(body) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL, err := useCase.CreateShortURL(string(body))
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	fullURL := ""

	if config.BaseURL == "" {
		var scheme string
		if req.TLS != nil {
			scheme = "https://"
		} else {
			scheme = "http://"
		}

		fullURL = scheme + req.Host + "/" + shortURL
	} else {
		fullURL = config.BaseURL + "/" + shortURL
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(fullURL))
}

func HandlerGet(res http.ResponseWriter, req *http.Request, useCase usecase.ShortURLUseCase) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(req, "shortURL")
	originalURL, err := useCase.GetOriginalURL(id)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(originalURL) > 4 {
		res.Header().Set("Location", originalURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}

func HandlerAPIPost(w http.ResponseWriter, r *http.Request, useCase usecase.ShortURLUseCase) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.Request
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	shortURL, err := useCase.CreateShortURL(req.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fullURL := ""

	if config.BaseURL == "" {
		var scheme string
		if r.TLS != nil {
			scheme = "https://"
		} else {
			scheme = "http://"
		}

		fullURL = scheme + r.Host + "/" + shortURL
	} else {
		fullURL = config.BaseURL + "/" + shortURL
	}

	resp := models.Response{
		Result: fullURL,
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
