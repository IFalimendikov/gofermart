package main

import (
	"context"
	"github.com/go-resty/resty/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gofermart/internal/config"
	"gofermart/internal/flag"
	"gofermart/internal/handler"
	"gofermart/internal/logger"
	"gofermart/internal/service"
	"gofermart/internal/storage"
	"gofermart/internal/transport"
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

	client := resty.New()

	s, err := service.New(log, cfg, store, client)
	if err != nil {
		log.Error("Error creating new gofermart service", "error", err)
	}
	go s.UpdateOrders(ctx)

	h := handler.New(s, log, &cfg)

	t := transport.New(&cfg, h, log)
	r := t.NewRouter()
	r.Run(cfg.ServerAddr)
}
