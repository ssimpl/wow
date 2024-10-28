package main

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type config struct {
	ServerAddr        string        `env:"SERVER_ADDR" env-default:":8080"`
	Requests          int           `env:"REQUESTS" env-default:"1"`
	SolutionTimeout   time.Duration `env:"SOLUTION_TIMEOUT" env-default:"5s"`
	ConnectionTimeout time.Duration `env:"CONNECTION_TIMEOUT" env-default:"5s"`
}

func newConfig() (config, error) {
	var cfg config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return cfg, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
