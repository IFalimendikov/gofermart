package main 

import (
	"context"
	"gofermart/internal/config"
	"gofermart/internal/flag"
	"gofermart/internal/logger"
	"gofermart/internal/service"
	"gofermart/internal/storage"
	"gofermart/internal/transport"
	"gofermart/internal/handler"
	"github.com/go-resty/resty/v2"
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
	// defer store.Drop(ctx)
	defer store.DB.Close()

	client := resty.New()

	s, err := service.New(log, cfg, store, client)
	if err != nil {
		log.Error("Error creating new gofermart service", "error", err)
	}
	go s.UpdateOrders(ctx)

	h := handler.New(s, log)

	t := transport.New(&cfg, h, log)
	r := t.NewRouter()
	r.Run(cfg.ServerAddr)
}