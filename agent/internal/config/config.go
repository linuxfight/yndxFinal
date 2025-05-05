package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

const (
	RequestTimeout = 2 * time.Second
)

type Config struct {
	ApiAddr        string
	ComputingPower int
}

func New() *Config {
	apiUrl := os.Getenv("API_ADDR")
	if apiUrl == "" {
		apiUrl = "localhost:9090"
	}

	powerStr := os.Getenv("POWER")
	power := 1
	if powerStr != "" {
		var err error
		if power, err = strconv.Atoi(powerStr); err != nil {
			log.Panicf("Error converting POWER to int: %s\n", err)
		}
	}

	return &Config{
		ApiAddr:        apiUrl,
		ComputingPower: power,
	}
}
