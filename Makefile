include ./configs/.config.env

.PHONY: build
build:
	go build -o api ./cmd/main.go

.PHONY: run
run:
	go run ./cmd/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: migrate_up
migrate_up:
	migrate -path ./internal/store/migrations -database $(DB_URL) -verbose up

.PHONY: migrate_down
migrate_down:
	migrate -path ./internal/store/migrations -database $(DB_URL) -verbose down


