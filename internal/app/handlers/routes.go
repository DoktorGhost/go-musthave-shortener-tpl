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

	// Общие middleware
	r.Use(logger.WithLogging)
	r.Use(compressor.GzipAndDecompressMiddleware)

	// Маршруты, к которым применяется UserMiddleware
	r.Group(func(authGroup chi.Router) {
		authGroup.Use(auth.UserMiddleware)

		authGroup.Post("/", func(w http.ResponseWriter, r *http.Request) {
			HandlerPost(w, r, useCase, conf)
		})
		authGroup.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
			HandlerAPIPost(w, r, useCase, conf)
		})
		authGroup.Post("/api/shorten/batch", func(w http.ResponseWriter, r *http.Request) {
			HandlerBatch(w, r, useCase, conf)
		})

	})

	// Маршруты, к которым не применяется UserMiddleware
	r.Get("/{shortURL}", func(w http.ResponseWriter, r *http.Request) {
		HandlerGet(w, r, useCase)
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		HandlerPing(w, r, conf)
	})
	r.Get("/api/user/urls", func(w http.ResponseWriter, r *http.Request) {
		HandlerGetUserURL(w, r, useCase)
	})

	return r
}
