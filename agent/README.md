# Агент

## Конфигурация
Происходит через изменение переменных среды:
```go
type Config struct {
	ApiAddr        string `env:"API_ADDR" env-default:"localhost:9090"` // адрес gRPC API
	ComputingPower int    `env:"COMPUTING_POWER" env-default:"10"` // количество воркеров
}
```

## Библиотеки
- cleanenv - для конфигурации
- grpc - для общения с оркестратором

## Структура проекта
```shell
.
├── Dockerfile
├── README.md
├── cmd
│ └── main.go
├── dev-compose.yml
├── go.mod
├── go.sum
├── internal
│ ├── app (приложение)
│ │ ├── app.go
│ │ ├── app_test.go (e2e тест)
│ │ ├── app_utils_test.go
│ │ └── utils.go
│ ├── config (конфиг приложения)
│ │ └── config.go
│ ├── tasks (взаимодействие с оркестратором)
│ │ ├── client.go
│ │ ├── gen
│ │ │ ├── tasks.pb.go
│ │ │ └── tasks_grpc.pb.go
│ │ ├── middlewares.go
│ │ ├── middlewares_test.go (unit тест)
│ │ └── middlewares_utils_test.go
│ └── worker (воркер и его методы)
│     └── worker.go
└── tasks.proto (grpc proto)
```
