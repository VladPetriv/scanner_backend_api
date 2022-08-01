# scanner_backend_api

## Description

scanner_backend_api is a backend side of tg_scanner application

Application will always create a dir:

- logs - for log's file

## Technology

Go, PostgreSQL, Gorilla-Mux, Zap, Go-Sqlmock, Testify, Golang-Migrate, JWT, Swagger

## Installation

```bash
 git clone git@github.com:VladPetriv/scanner_backend_api.git

 cd scanner_backend_api

 go mod download

```

## Before start

Please create .env file with this fields:

- DATABASE_URL = PostgreSQL user
- LOG_LEVEL = Level which logger will handle
- MIGRATIONS_PATH = Path to migrations:"file://./db/migrations"
- PORT = Bind address which server going to use
- JWT_SECRET_KEY = Secret key for json web token

## Usage

Start with docker-compose:

```bash
 make docker
```

Start an application locally:

```bash
 make run # Or you can use go run ./cmd/main.go
```

Running test suite:

```bash
 make mock # Run it if mocks folder is not exist

 make test
```
