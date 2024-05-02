package handlers

import (
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

	//проверка реальности url
	resp, err := http.Get(string(body))
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL, err := useCase.CreateShortURL(string(body))
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var scheme string
	if req.TLS != nil {
		scheme = "https://"
	} else {
		scheme = "http://"
	}

	fullURL := scheme + req.Host + "/" + shortURL

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
