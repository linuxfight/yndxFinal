# Агент

# Разработка
1. Установите NodeJS (https://nodejs.org/)
```shell
go run cmd/main.go
```

# Структура проекта
```shell
.
├── Dockerfile
├── README.md (этот файл)
├── cmd
│ └── main.go
├── dev-compose.yml
├── go.mod
├── go.sum
├── internal
│ ├── config
│ │ └── config.go
│ ├── tasks
│ │ ├── client.go
│ │ └── gen
│ │     ├── tasks.pb.go
│ │     └── tasks_grpc.pb.go
│ └── worker
│     └── worker.go
└── tasks.proto
```
