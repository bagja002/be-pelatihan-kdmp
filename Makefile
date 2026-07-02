.PHONY: run gen tidy build keygen

## keygen: print fresh secrets for .env (JWT_SECRET + ENCRYPTION_KEY)
keygen:
	@echo "JWT_SECRET=$$(openssl rand -base64 48)"
	@echo "ENCRYPTION_KEY=$$(openssl rand -hex 32)"

## run: start the API server
run:
	go run ./cmd/api

## gen: scaffold a new entity (usage: make gen name=Product)
gen:
	@if [ -z "$(name)" ]; then echo "usage: make gen name=Product"; exit 1; fi
	go run ./cmd/generate -name $(name)

## tidy: resolve go module dependencies
tidy:
	go mod tidy

## build: compile the API binary into ./bin/api
build:
	go build -o bin/api ./cmd/api
