# Финальный проект 📟

## Онлайн демо 🌐
Если вам не подходит ни один из способов запуска, то можете воспользоваться [онлайн демо](https://calc.lxft.tech/).

## Структура проекта 📁
В папках orchestrator, agent и frontend есть описание решений, библиотек и файлов.

## Документация API 📃
Документация запросов и ответом доступна по [этому адресу](https://calc-backend.lxft.tech/docs#/).
Здесь можно отправить запросы на бэкэнд, посмотреть примерные данные, а также пути и то, что делают запросы.

## Функционал ✅
- CI - тестирование, сборка Docker образов (Готово)
- Решение простых математических выражений из +,-,*,/,(,) и чисел (Готово)
- Веб интерфейс (Готово)
- Документация (Готово)
- Решение выражений из !,^,% (В разработке)
- OpenTelemetry (Logging, Tracing, Metrics) (В разработке)
- Выбор Scope при получении выражений (Свои/Все) (В разработке)
- CD - автоматическое обновление Docker образов на сервере с помощью WatchTower (В разработке)

## Как это работает? 🧪
![explain](./content/explain.jpg "Объяснение")

## Схема СУБД
1. Valkey
```shell
taskId -> taskId;arg1;arg2;op;res
```
2. PostgreSQL
![db](./content/db.png "Схема СУБД")

## Запуск 🚀
### 1. Docker
1. [Установите Docker](https://www.docker.com/products/docker-desktop/)
2. Откройте папку с проектом в терминале
3. Пропишите:
```shell
docker compose up
```
### 2. Aeza (нужно, если нет возможности установить Docker)
1. [Перейдите на Terminator](https://terminator.aeza.net/ru/)
2. Следуйте шагам из части 1, но чтобы скачать проект - скачайте [zip архив с GitHub](https://github.com/linuxfight/yndxFinal/archive/refs/heads/main.zip)
### 3. В ручную
1. [Установите PostgreSQL](https://www.pgadmin.org/)
2. Создайте бд со следующими параметрами [(см документацию оркестратора)](./orchestrator/README.md)
3. [Установите Valkey](https://valkey.io/topics/installation/)
4. [Установите Go](https://go.dev/doc/install)
5. В терминале из папки проекта запустите в разных окнах:
```shell
# запуск оркестратора
cd orchestrator
go run cmd/main.go

# запуск агента
cd agent
go run cmd/main.go
```

## Тестирование 🛠
Тестирование завязано на TestContainers, поэтому Docker обязателен, смотрите README в папках agent и orchestrator для информации о каждом тесте.
```shell
# agent
cd agent && go test -v -cover ./...

# orchestrator
cd orchestrator && go test -v -cover ./...
```


## Фидбэк 🖋
Если вам не трудно, то напишите ваш Фидбэк по решению в issues :)
