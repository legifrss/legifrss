.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


build: ## build
	@go build ./server.go

test:  build ## test the app
	@go test ./...

run: test ## run the app
	@go run ./server.go

run-commit: run ## run the app and commit feeds
	@scripts/commit.sh

delete-feeds: ## delete all feeds
	@rm -rf feed/*
