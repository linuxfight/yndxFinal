package config

import (
	"os"
	"strconv"
)

type Config struct {
	ValkeyConn   string
	PostgresConn string
	JwtSecret    string

	AddictionTime      int
	SubstractionTime   int
	MultiplicationTime int
	DivisionTime       int

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

	// TODO: TIME_ADDITION_MS
	// TODO: TIME_SUBTRACTION_MS
	// TODO: TIME_MULTIPLICATIONS_MS
	// TODO: TIME_DIVISIONS_MS
	// TODO: port to clearenv

	return &Config{
		ValkeyConn:    valkeyConn,
		PostgresConn:  postgresConn,
		OperationTime: operationTime,
	}
}
