package server

import (
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/handlers"
	"net/http"
)

func StartServer(port string) error {

	mux := handlers.InitMux()
	handlers.InitRoutes(mux)

	err := http.ListenAndServe(":"+port, mux)

	if err != nil {
		return err
	}
	return nil
}
