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
	cfgN := Config{}
	err := env.Load()
	if err != nil {
		return err
	}

	cfgN.ServerAddr = os.Getenv("RUN_ADDRESS")
	if cfgN.ServerAddr != "" {
		cfg.ServerAddr = cfgN.ServerAddr
	}
	
	cfg.DatabaseURI = os.Getenv("DATABASE_URI")
	if cfgN.DatabaseURI != "" {
		cfg.DatabaseURI = cfgN.DatabaseURI
	}

	cfg.AccrualAddr = os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	if cfgN.AccrualAddr != "" {
		cfg.AccrualAddr = cfgN.AccrualAddr
	}

	return nil
}