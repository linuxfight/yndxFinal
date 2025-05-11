package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"time"
)

const (
	RequestTimeout = 1 * time.Second
)

type Config struct {
	ApiAddr        string `env:"API_ADDR" env-default:"localhost:9090"`
	ComputingPower int    `env:"COMPUTING_POWER" env-default:"10"`
}

func New() *Config {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Panicf("failed to read environment variables: %v", err)
	}
	return &cfg
}
