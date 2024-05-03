package server

import (
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/storage/maps"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/usecase"
	"net/http"
	"strconv"
)

func StartServer(port int) error {
	db := maps.NewMapStorage()
	shortURLUseCase := usecase.NewShortURLUseCase(db)

	r := handlers.InitRoutes(*shortURLUseCase)

	err := http.ListenAndServe(":"+strconv.Itoa(port), r)

	if err != nil {
		return err
	}
	return nil
}
