package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"gofermart/internal/models"
	"gofermart/internal/storage"
)

type Service interface {

}

type Gofermart struct {
	Log *slog.Logger
	Storage *storage.Storage
}

func New(log *slog.Logger, storage *storage.Storage) (*Gofermart, error){
	service := Gofermart {
		Log: log,
		Storage: storage,
	}
	return &service, nil
}