package server

import (
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/storage/maps"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/usecase"
	"net/http"
)

func StartServer(port string) error {
	db := maps.NewMapStorage()
	shortUrlUseCase := usecase.NewShortUrlUseCase(db)

	mux := handlers.InitMux()
	handlers.InitRoutes(mux, *shortUrlUseCase)

	err := http.ListenAndServe(":"+port, mux)

	if err != nil {
		return err
	}
	return nil
}
