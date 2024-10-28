package main

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type config struct {
	Addr                 string        `env:"LISTEN_ADDR" env-default:":8080"`
	Difficulty           int           `env:"POW_DIFFICULTY" env-default:"4"`
	ClientWaitingTimeout time.Duration `env:"CLIENT_WAITING_TIMEOUT" env-default:"5s"`
	ShutdownTimeout      time.Duration `env:"SHUTDOWN_TIMEOUT" env-default:"5s"`
}

func newConfig() (config, error) {
	var cfg config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return cfg, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
