.PHONY: build run dev clean tidy lint test

APP_NAME=slack-link-bot
MAIN_PATH=cmd/server/main.go

build:
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

dev:
	set -a && source config/.env && set +a && go run $(MAIN_PATH)

run: build
	set -a && source config/.env && set +a && ./bin/$(APP_NAME)

clean:
	rm -rf bin/

tidy:
	go mod tidy

lint:
	golangci-lint run

test:
	set -a && source .env && set +a && go test -v ./...

.DEFAULT_GOAL := dev