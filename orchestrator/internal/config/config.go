package config

import (
	"os"
	"strconv"
)

type Config struct {
	ValkeyConn   string
	PostgresConn string

	OperationTime int
}

func New() *Config {
	valkeyConn := os.Getenv("VALKEY_CONN")
	if valkeyConn == "" {
		valkeyConn = "127.0.0.1:6379"
	}

	postgresConn := os.Getenv("POSTGRES_CONN")
	if postgresConn == "" {
		postgresConn = "postgres://postgres:password@localhost:5432/db"
	}

	operationTime := 1000
	var err error
	operationTimeStr := os.Getenv("OPERATION_TIME")
	if operationTimeStr != "" {
		if operationTime, err = strconv.Atoi(operationTimeStr); err != nil {
			panic(err)
		}
	}

	return &Config{
		ValkeyConn:    valkeyConn,
		PostgresConn:  postgresConn,
		OperationTime: operationTime,
	}
}
