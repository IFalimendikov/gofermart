package handler

import (
	"context"
	"log/slog"
	"gofermart/internal/models"
)

type Service interface {

}

type Handler struct {
	Service Service
	log *slog.Logger
}

func New(s Service, log *slog.Logger) *Handler {
	return &Handler{
		Service: s,
		log: log,
	}
}