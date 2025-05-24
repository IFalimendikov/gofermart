package flag

import (
	"flag"

	"gofermart/internal/config"
)

func Parse() config.Config {
	cfg := config.Config{}

	flag.StringVar(&cfg.ServerAddr, "a", cfg.ServerAddr, "Address and Port")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "Database URL")
	flag.StringVar(&cfg.AccrualAddr, "r", cfg.AccrualAddr, "Accrual Service address")
	flag.Parse()

	return cfg
}