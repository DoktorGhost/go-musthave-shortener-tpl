package main

import (
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/server"
)

func main() {

	err := server.StartServer()
	if err != nil {
		panic(err)
	}
}
