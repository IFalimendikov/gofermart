package main 

import (
	"context"
	"gofermart/internal/config"
	"gofermart/internal/flag"
	"gofermart/internal/logger"
	"gofermart/internal/services"
	"gofermart/internal/storage"
	"gofermart/internal/transport"
	"gofermart/internal/handler"
)

func main() {
	log := logger.New()

	cfg := flag.Parse()
	err := config.Read(&cfg)
	if err != nil {
		log.Error("Error reading env file", "error", err)
	}

	ctx := context.Background()

	store, err := storage.New(ctx, &cfg)
	if err != nil {
		log.Error("Error creating new storage", "error", err)
	}
	defer store.DB.Close()

	s := services.NewGofermartService(ctx, log, store)
	h := handler.NewHandler(s, log)

	t := transport.NewTransport(cfg, h, log)
	r := transport.NewRouter(t)
	r.Run(cfg.ServerAddr)
}