package config

import (
	"os"

	env "github.com/joho/godotenv"
)

type Config struct {
    ServerAddr  string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
    DatabaseURI string `env:"DATABASE_URI"`
    AccrualAddr string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}
func Read(cfg *Config) error {
    err := env.Load()
    if err != nil && !os.IsNotExist(err) {
        return err
    }

    if cfg.ServerAddr == "" {
        cfg.ServerAddr = os.Getenv("RUN_ADDRESS")
    }

    if cfg.DatabaseURI == "" {
        cfg.DatabaseURI = os.Getenv("DATABASE_URI")
    }

    if cfg.AccrualAddr == "" {
        cfg.AccrualAddr = os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
    }

    return nil
}