.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


build: ## build
	@go build ./server.go

test:  build ## test the app
	@go test ./...

run: test ## run the app
	@go run ./server.go

run-batch: run
	@scripts/deploy.sh

delete-feeds:
	@rm -rf feed/*
