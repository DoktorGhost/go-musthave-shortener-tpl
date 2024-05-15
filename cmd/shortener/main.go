package main

import (
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/config"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/server"
)

func main() {
	conf := config.ParseConfig()

	err := server.StartServer(conf)
	if err != nil {
		panic(err)
	}
}
