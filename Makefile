.PHONY: build run test docker-up docker-down migrate

build:
	go build -o bin/user-sync-service ./cmd/api

run:
	go run ./cmd/api

test:
	go test -v -race -cover ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate:
	migrate -path ./scripts/migrations -database "mysql://root:root@tcp(localhost:3308)/user_sync" up

migrate-down:
	migrate -path ./scripts/migrations -database "mysql://root:root@tcp(localhost:3308)/user_sync" down

lint:
	golangci-lint run

tidy:
	go mod tidy