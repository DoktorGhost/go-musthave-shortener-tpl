package handlers

import (
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

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		HandlerPost(w, r, useCase, conf)
	})
	r.Get("/{shortURL}", func(w http.ResponseWriter, r *http.Request) {
		HandlerGet(w, r, useCase)
	})
	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		HandlerAPIPost(w, r, useCase, conf)
	})
	return r
}
