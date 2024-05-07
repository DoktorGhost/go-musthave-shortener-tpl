package main

import (
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/config"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/server"
)

func main() {
	hostPort := config.ParseConfig()

	err := server.StartServer(hostPort)
	if err != nil {
		panic(err)
	}
}
