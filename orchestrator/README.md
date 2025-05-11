```shell
.
├── Dockerfile
├── README.md
├── cmd
│ └── main.go
├── dev-compose.yml
├── docs
│ ├── docs.go
│ ├── swagger.json
│ └── swagger.yaml
├── go.mod
├── go.sum
├── internal
│ ├── app
│ │ ├── app.go
│ │ ├── app_grpc_test.go
│ │ ├── app_http_test.go
│ │ └── app_utils_test.go
│ ├── config
│ │ └── config.go
│ ├── controllers
│ │ ├── auth
│ │ │ ├── controller.go
│ │ │ ├── login.go
│ │ │ ├── password.go
│ │ │ ├── password_test.go
│ │ │ ├── register.go
│ │ │ ├── utils.go
│ │ │ └── utils_test.go
│ │ ├── dto
│ │ │ ├── auth.go
│ │ │ ├── const.go
│ │ │ ├── error.go
│ │ │ ├── solverCalculate.go
│ │ │ └── solverGet.go
│ │ ├── expr
│ │ │ ├── calculate.go
│ │ │ ├── controller.go
│ │ │ ├── getById.go
│ │ │ └── list.go
│ │ ├── fiber.go
│ │ ├── middlewares
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
│ │ └── tasks
│ │     ├── gen
│ │     │ ├── tasks.pb.go
│ │     │ └── tasks_grpc.pb.go
│ │     ├── server.go
│ │     └── utils.go
│ ├── db
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
│ └── utils
│     ├── fiber.go
│     └── fiber_test.go
├── pkg
│ └── calc
│     ├── calc.go
│     └── calc_test.go
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