# Оркестратор

## Конфигурация
Происходит через изменение переменных среды:
```go
type Config struct {
	ValkeyConn   string `env:"VALKEY_CONN" env-default:"127.0.0.1:6379"` // подключение к Redis/Valkey
	PostgresConn string `env:"POSTGRES_CONN" env-default:"postgres://postgres:password@localhost:5432/db"` // подключение к PostgreSQL
	JwtSecret    string `env:"JWT_SECRET" env-default:"not_v3ry_s3cR3T"` // секрет для генерации и проверки JWT токенов

	AddictionTime      int `env:"TIME_ADDITION_MS" env-default:"1000"` // время на сложение
	SubstractionTime   int `env:"TIME_SUBTRACTION_MS" env-default:"1000"` // время на вычитание
	MultiplicationTime int `env:"TIME_MULTIPLICATIONS_MS" env-default:"1000"` // время на умножение
	DivisionTime       int `env:"TIME_DIVISIONS_MS" env-default:"1000"` // время на деление
}
```

## Библиотеки
- argon2 - хэширование паролей
- sonic - для работы с json
- testcontainers - e2e тестирование
- grpc - общение с агентом
- valkey, pgx - общение с бд
- sqlc - генерация кода для бд
- clearenv - конфигурация
- ulid - удобные идентификаторы
- swag - генерация swagger
- fiber - http фреймворк

## Структура проекта
```shell
.
├── Dockerfile
├── README.md
├── cmd
│ └── main.go
├── dev-compose.yml
├── docs (документация swagger)
│ ├── docs.go
│ ├── swagger.json
│ └── swagger.yaml
├── go.mod
├── go.sum
├── internal
│ ├── app (приложение)
│ │ ├── app.go
│ │ ├── app_grpc_test.go (e2e тест)
│ │ ├── app_http_test.go (e2e тест)
│ │ └── app_utils_test.go
│ ├── config (конфиг)
│ │ └── config.go
│ ├── controllers (http и grpc логика)
│ │ ├── auth (login, register)
│ │ │ ├── controller.go
│ │ │ ├── login.go
│ │ │ ├── password.go
│ │ │ ├── password_test.go
│ │ │ ├── register.go
│ │ │ ├── utils.go
│ │ │ └── utils_test.go
│ │ ├── dto (модели http API)
│ │ │ ├── auth.go
│ │ │ ├── const.go
│ │ │ ├── error.go
│ │ │ ├── solverCalculate.go
│ │ │ └── solverGet.go
│ │ ├── expr (calculate, list, id)
│ │ │ ├── calculate.go
│ │ │ ├── controller.go
│ │ │ ├── getById.go
│ │ │ └── list.go
│ │ ├── fiber.go
│ │ ├── middlewares (мидлвари для http api)
│ │ │ ├── jwt
│ │ │ │ ├── config.go
│ │ │ │ ├── config_test.go
│ │ │ │ ├── crypto.go
│ │ │ │ ├── jwt.go
│ │ │ │ ├── jwt_test.go
│ │ │ │ └── utils.go
│ │ │ ├── recoverer
│ │ │ │ ├── recoverer.go
│ │ │ │ └── recoverer_test.go
│ │ │ └── swagger
│ │ │     ├── swagger.go
│ │ │     ├── swagger.json
│ │ │     ├── swagger.yaml
│ │ │     ├── swagger_missing.json
│ │ │     └── swagger_test.go
│ │ └── tasks (grpc api)
│ │     ├── gen
│ │     │ ├── tasks.pb.go
│ │     │ └── tasks_grpc.pb.go
│ │     ├── server.go
│ │     └── utils.go
│ ├── db (codegen и код для valkey, postgres)
│ │ ├── cache.go
│ │ ├── cache_test.go
│ │ ├── connection.go
│ │ ├── expressions
│ │ │ ├── db.go
│ │ │ ├── models.go
│ │ │ └── query.sql.go
│ │ └── users
│ │     ├── db.go
│ │     ├── models.go
│ │     └── query.sql.go
│ └── utils (утилиты для http api)
│     ├── fiber.go
│     └── fiber_test.go (unit тест)
├── pkg
│ └── calc (парсер)
│     ├── calc.go
│     └── calc_test.go (unit тест)
├── sql
│ ├── expressions
│ │ ├── query.sql
│ │ └── schema.sql
│ └── users
│     ├── query.sql
│     └── schema.sql
├── sqlc.yml
└── tasks.proto
```
