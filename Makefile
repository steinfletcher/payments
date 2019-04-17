.PHONY: build test postgres-dev test-ingegration test-all run

DB_ADDR = 'postgresql://localhost:15432/payments?user=postgres&password=postgres&sslmode=disable'

build:
	go build -o payments cmd/payments/main.go

test:
	go test -v ./...

test-integration:
	DB_ADDR=$(DB_ADDR) go test -v ./integration
	DB_ADDR=$(DB_ADDR) go test -v ./postgres

test-all: test test-integration

run:
	DB_ADDR=$(DB_ADDR) go run cmd/payments/main.go

postgres-dev:
	@docker stop payments || true && docker rm payments || true
	docker run -d -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=payments -p 15432:5432 --name payments postgres:9.6
