package handlers

import (
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/usecase"
	"io"
	"net/http"
)

func HandlerPost(res http.ResponseWriter, req *http.Request, useCase usecase.ShortUrlUseCase) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL := useCase.CreateShortUrl(string(body))

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

func HandlerGet(res http.ResponseWriter, req *http.Request, useCase usecase.ShortUrlUseCase) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := req.URL.Path[1:]
	originalURL, err := useCase.GetOriginalUrl(id)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if originalURL != "" {
		res.Header().Set("Location", originalURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}
