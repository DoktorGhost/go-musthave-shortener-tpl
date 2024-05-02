package handlers

import (
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/usecase"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func InitRoutes(useCase usecase.ShortUrlUseCase) chi.Router {
	r := chi.NewRouter()

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		HandlerPost(w, r, useCase)
	})
	r.Get("/{shortURL}", func(w http.ResponseWriter, r *http.Request) {
		HandlerGet(w, r, useCase)
	})
	return r
}
