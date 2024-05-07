package handlers

import (
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/logger"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/usecase"
	"github.com/go-chi/chi/v5"
	"net/http"
)

/*
func InitRoutes(useCase usecase.ShortURLUseCase) chi.Router {
	r := chi.NewRouter()

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		HandlerPost(w, r, useCase)
	})
	r.Get("/{shortURL}", func(w http.ResponseWriter, r *http.Request) {
		HandlerGet(w, r, useCase)
	})
	return r
}


*/

func InitRoutes(useCase usecase.ShortURLUseCase) chi.Router {
	r := chi.NewRouter()

	r.Use(logger.WithLogging)

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		HandlerPost(w, r, useCase)
	})
	r.Get("/{shortURL}", func(w http.ResponseWriter, r *http.Request) {
		HandlerGet(w, r, useCase)
	})
	return r
}
