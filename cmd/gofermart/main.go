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
	cfg := flag.ParseFlags()
	config.Read(&cfg)

	log := logger.NewLogger()
	ctx := context.Background()

	store, err := storage.NewStorage(ctx, &cfg)
	if err != nil {
		log.Error("Error creating new storage", "error", err)
	}
	defer store.File.Close()
	defer store.DB.Close()

	s := services.NewGofermartService(ctx, log, store)
	h := handler.NewHandler(s, log)

	t := transport.NewTransport(cfg, h, log)
	r := transport.NewRouter(t)
	r.Run(cfg.ServerAddr)
}