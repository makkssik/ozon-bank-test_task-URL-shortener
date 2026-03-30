.PHONY: run-local-memory run-local-postgres docker-memory docker-postgres test

run-local-memory:
	go run ./cmd/url-shortener/main.go -storage=memory

run-local-postgres:
	docker-compose up -d db migrator
	go run ./cmd/url-shortener/main.go -storage=postgres

docker-memory:
	docker-compose -f docker-compose.memory.yml up --build

docker-postgres:
	docker-compose up --build

test:
	go test -race -v ./tests/...