package server

import (
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/storage/maps"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/usecase"
	"net/http"
)

func StartServer(port string) error {
	db := maps.NewMapStorage()
	shortUrlUseCase := usecase.NewShortURLUseCase(db)

	r := handlers.InitRoutes(*shortUrlUseCase)

	err := http.ListenAndServe(":"+port, r)

	if err != nil {
		return err
	}
	return nil
}
