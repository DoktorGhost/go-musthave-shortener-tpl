package handlers

import (
	"net/http"
)

func InitMux() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}

func InitRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", HandlerPost)
	mux.HandleFunc("/{ip}", HandlerGet)
}
