package main

import (
	"github.com/saptaka/pos/config"
	"github.com/saptaka/pos/repository"
	"github.com/saptaka/pos/server"
)

func main() {
	cfg := config.Setup()
	repo := repository.NewRepository(cfg)
	service := server.NewServer(cfg, repo)
	service.Listen(3030)
}
