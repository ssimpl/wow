.PHONY: test
test:
	go test -timeout 1m -race -cover ./...

.PHONY: lint
lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.61 golangci-lint run -v

.PHONY: up
up:
	docker-compose up --build || docker compose up --build

.PHONY: down
down:
	docker-compose down || docker compose down
