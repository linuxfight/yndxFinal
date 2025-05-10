package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"orchestrator/internal/controllers/tasks/gen"
)

type Config struct {
	ValkeyConn   string `env:"VALKEY_CONN" env-default:"127.0.0.1:6379"`
	PostgresConn string `env:"POSTGRES_CONN" env-default:"postgres://postgres:password@localhost:5432/db"`
	JwtSecret    string `env:"JWT_SECRET" env-default:"not_v3ry_s3cR3T"`

	AddictionTime      int `env:"TIME_ADDITION_MS" env-default:"1000"`
	SubstractionTime   int `env:"TIME_SUBTRACTION_MS" env-default:"1000"`
	MultiplicationTime int `env:"TIME_MULTIPLICATIONS_MS" env-default:"1000"`
	DivisionTime       int `env:"TIME_DIVISIONS_MS" env-default:"1000"`
}

func (cfg *Config) GetOperationTime(op gen.Operator) int32 {
	switch op {
	case gen.Operator_ADDICTION:
		return int32(cfg.AddictionTime)
	case gen.Operator_SUBTRACTION:
		return int32(cfg.SubstractionTime)
	case gen.Operator_MULTIPLICATION:
		return int32(cfg.MultiplicationTime)
	case gen.Operator_DIVISION:
		return int32(cfg.DivisionTime)
	}

	log.Panicf("invalid op: %s", op)
	return 0
}

func New() *Config {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Panicf("failed to read environment variables: %v", err)
	}
	return &cfg
}
