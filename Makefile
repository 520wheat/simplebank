ifneq (,$(wildcard .env))
    include .env
    export
endif

GOPATH := $(shell go env GOPATH 2>/dev/null || echo $(HOME)/go)
PATH := $(GOPATH)/bin:$(PATH)

DB_URL ?= postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker compose up -d postgres

postgres-down:
	docker compose down

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres postgres-down migrateup migratedown sqlc test server
