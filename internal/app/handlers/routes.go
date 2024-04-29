package handlers

import (
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/usecase"
	"net/http"
)

func InitMux() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}

func InitRoutes(mux *http.ServeMux, useCase usecase.ShortUrlUseCase) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HandlerPost(w, r, useCase)
	})

	mux.HandleFunc("/{shortURL}", func(w http.ResponseWriter, r *http.Request) {
		HandlerGet(w, r, useCase)
	})

}
