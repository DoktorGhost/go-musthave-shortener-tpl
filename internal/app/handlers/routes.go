package handlers

import (
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/usecase"
	"github.com/gorilla/mux"
	"net/http"
)

func InitMux() *mux.Router {
	router := mux.NewRouter()
	return router
}

func InitRoutes(mux *mux.Router, useCase usecase.ShortUrlUseCase) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HandlerPost(w, r, useCase)
	}).Methods("POST")

	mux.HandleFunc("/{shortURL}", func(w http.ResponseWriter, r *http.Request) {
		HandlerGet(w, r, useCase)
	}).Methods("GET")
}
