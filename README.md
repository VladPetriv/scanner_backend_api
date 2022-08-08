
# scanner_backend_api

scanner_backend_api is a backend side represented as REST-API for processing data from tg_scanner application.

## Tech Stack

**Server:** 
- Gorilla-Mux
- JWT
- Apache Kafka

**DB:**
- PostgreSQL
- Golang-Migrate

**Testing:**
- Testify
- Go-Sqlmock

**Docs:**
- Swagger


## Features

- Microservice communication with Apache Kafka
- JWT authentication
- Unit test for all application layers [handler -> service -> repository]
- Generated swagger documentation


## Environment Variables

To run this project, you will need to add the following environment variables to your ".config.env" file which locate in "config" folder:

- `DATABASE_URL` - PostgreSQL url
- `LOG_LEVEL` - Level which logger will process
- `MIGRATIONS_PATH` - Path to migrations: "file://./db/migrations/"
- `PORT` - Bind address which server will use
- `JWT_SECRET_KEY` - Secret key for Json Web Token
- `KAFKA_ADDR` - Apache Kafka address

## Run Locally

Clone the project

```bash
  git clone git@github.com:VladPetriv/scanner_backend_api.git
```

Go to the project directory

```bash
  cd scanner_backend_api
```

Install dependencies

```bash
  go mod download
```

Start the application:

```bash
  # Make sure that Apache Kafka and PostgreSQL are running
  make run # Or you can use "go run ./cmd/main.go"
```

Start the application with docker-compose:

```bash
  make docker
```
## Running Tests

To run tests, run the following command:

```bash
  # Run it only if "mocks" folder not exist or if you updated "service.go" or "repository.go" files
  make mock 
```

```bash
  make test 
```

## License

[MIT](https://choosealicense.com/licenses/mit/)
