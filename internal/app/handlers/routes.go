package handlers

import (
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/auth"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/compressor"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/config"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/logger"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/usecase"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func InitRoutes(useCase usecase.ShortURLUseCase, conf *config.Config) chi.Router {
	r := chi.NewRouter()

	r.Use(logger.WithLogging)
	r.Use(compressor.GzipAndDecompressMiddleware)

	r.Get("/{shortURL}", func(w http.ResponseWriter, r *http.Request) {
		HandlerGet(w, r, useCase)
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		HandlerPing(w, r, conf)
	})
	r.Get("/api/user/urls", func(w http.ResponseWriter, r *http.Request) {
		HandlerGetUserURL(w, r, useCase)
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.UserMiddleware)
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			HandlerPost(w, r, useCase, conf)
		})

		r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
			HandlerAPIPost(w, r, useCase, conf)
		})
		r.Post("/api/shorten/batch", func(w http.ResponseWriter, r *http.Request) {
			HandlerBatch(w, r, useCase, conf)
		})
	})
	return r
}
